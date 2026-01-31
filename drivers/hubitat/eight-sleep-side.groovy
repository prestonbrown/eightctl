/**
 * Eight Sleep Side
 * Child device driver for one side of Eight Sleep Pod
 * Communicates with eightctl hubitat server via parent device
 *
 * Copyright 2024 eightctl
 * Licensed under the Apache License, Version 2.0
 */

metadata {
    definition(name: "Eight Sleep Side", namespace: "eightctl", author: "eightctl") {
        capability "Switch"
        capability "Thermostat"
        capability "TemperatureMeasurement"
        capability "Refresh"

        // Eight Sleep level: -100 (coolest) to +100 (warmest)
        attribute "level", "number"
        attribute "targetLevel", "number"
        attribute "isActive", "string"

        command "setLevel", [[name: "level", type: "NUMBER", description: "Temperature level (-100 to 100)"]]
        command "levelUp"
        command "levelDown"
    }

    preferences {
        input name: "levelStep", type: "number", title: "Level Step (for up/down)", defaultValue: 10, required: true
    }
}

def installed() {
    logDebug("Eight Sleep Side installed")
    initialize()
}

def updated() {
    logDebug("Eight Sleep Side updated")
    initialize()
}

def initialize() {
    logDebug("Initializing Eight Sleep Side")

    // Set initial thermostat mode
    sendEvent(name: "thermostatMode", value: "auto")
    sendEvent(name: "supportedThermostatModes", value: ["off", "auto"])
    sendEvent(name: "supportedThermostatFanModes", value: [])
}

// Switch capability
def on() {
    logDebug("Turning on ${getSide()} side")

    def params = [
        uri: "http://${getServerIP()}:${getServerPort()}",
        path: "/${getSide()}/on",
        contentType: "application/json",
        timeout: 10
    ]

    try {
        httpPut(params) { response ->
            if (response.status == 200 || response.status == 204) {
                logDebug("Successfully turned on ${getSide()} side")
                sendEvent(name: "switch", value: "on")
                sendEvent(name: "isActive", value: "true")
                sendEvent(name: "thermostatMode", value: "auto")
                parent.refresh()
            } else {
                log.error("Unexpected response turning on: ${response.status}")
            }
        }
    } catch (Exception e) {
        log.error("Error turning on ${getSide()} side: ${e.message}")
    }
}

def off() {
    logDebug("Turning off ${getSide()} side")

    def params = [
        uri: "http://${getServerIP()}:${getServerPort()}",
        path: "/${getSide()}/off",
        contentType: "application/json",
        timeout: 10
    ]

    try {
        httpPut(params) { response ->
            if (response.status == 200 || response.status == 204) {
                logDebug("Successfully turned off ${getSide()} side")
                sendEvent(name: "switch", value: "off")
                sendEvent(name: "isActive", value: "false")
                sendEvent(name: "thermostatMode", value: "off")
                parent.refresh()
            } else {
                log.error("Unexpected response turning off: ${response.status}")
            }
        }
    } catch (Exception e) {
        log.error("Error turning off ${getSide()} side: ${e.message}")
    }
}

// Level control
def setLevel(level) {
    def intLevel = level.toInteger()

    // Clamp to valid range
    if (intLevel < -100) intLevel = -100
    if (intLevel > 100) intLevel = 100

    logDebug("Setting ${getSide()} side level to ${intLevel}")

    def params = [
        uri: "http://${getServerIP()}:${getServerPort()}",
        path: "/${getSide()}/temperature",
        query: [level: intLevel.toString()],
        contentType: "application/json",
        timeout: 10
    ]

    try {
        httpPut(params) { response ->
            if (response.status == 200 || response.status == 204) {
                logDebug("Successfully set ${getSide()} side level to ${intLevel}")
                sendEvent(name: "level", value: intLevel)
                sendEvent(name: "targetLevel", value: intLevel)
                updateThermostatFromLevel(intLevel)
                parent.refresh()
            } else {
                log.error("Unexpected response setting level: ${response.status}")
            }
        }
    } catch (Exception e) {
        log.error("Error setting ${getSide()} side level: ${e.message}")
    }
}

def levelUp() {
    def currentLevel = device.currentValue("level") ?: 0
    def step = levelStep ?: 10
    setLevel(currentLevel + step)
}

def levelDown() {
    def currentLevel = device.currentValue("level") ?: 0
    def step = levelStep ?: 10
    setLevel(currentLevel - step)
}

// Thermostat capability
// Maps Eight Sleep level (-100 to 100) to temperature concept
// Level -100 = cooling to ~55F, Level 0 = neutral ~70F, Level 100 = heating to ~110F

def setHeatingSetpoint(temp) {
    logDebug("setHeatingSetpoint called with ${temp}")
    // Convert temperature to Eight Sleep level
    // Assume temp is in Fahrenheit: 55F = -100, 70F = 0, 110F = 100
    def level = mapTemperatureToLevel(temp)
    if (level > 0) {
        setLevel(level)
    } else {
        log.warn("Temperature ${temp}F maps to cooling level ${level}, use setCoolingSetpoint instead")
    }
}

def setCoolingSetpoint(temp) {
    logDebug("setCoolingSetpoint called with ${temp}")
    // Convert temperature to Eight Sleep level
    def level = mapTemperatureToLevel(temp)
    if (level < 0) {
        setLevel(level)
    } else {
        log.warn("Temperature ${temp}F maps to heating level ${level}, use setHeatingSetpoint instead")
    }
}

def setThermostatMode(mode) {
    logDebug("setThermostatMode called with ${mode}")
    if (mode == "off") {
        off()
    } else if (mode == "auto" || mode == "heat" || mode == "cool") {
        on()
        sendEvent(name: "thermostatMode", value: "auto")
    }
}

def setThermostatFanMode(mode) {
    logDebug("setThermostatFanMode called - not applicable for Eight Sleep")
}

// Not applicable but required for thermostat capability
def auto() { setThermostatMode("auto") }
def heat() { setThermostatMode("auto") }
def cool() { setThermostatMode("auto") }
def emergencyHeat() { log.warn("emergencyHeat not supported") }
def fanAuto() { log.debug("fanAuto not applicable") }
def fanCirculate() { log.debug("fanCirculate not applicable") }
def fanOn() { log.debug("fanOn not applicable") }

def refresh() {
    logDebug("Refreshing ${getSide()} side via parent")
    parent.refresh()
}

// Called by parent to update status
def updateStatus(Map data) {
    logDebug("Updating status with: ${data}")

    if (data.containsKey("currentLevel")) {
        sendEvent(name: "level", value: data.currentLevel)
        updateThermostatFromLevel(data.currentLevel)
    }

    if (data.containsKey("targetLevel")) {
        sendEvent(name: "targetLevel", value: data.targetLevel)
    }

    if (data.containsKey("isActive")) {
        def isActive = data.isActive
        sendEvent(name: "isActive", value: isActive.toString())
        sendEvent(name: "switch", value: isActive ? "on" : "off")
        sendEvent(name: "thermostatMode", value: isActive ? "auto" : "off")
    }

    if (data.containsKey("currentTemperature")) {
        // Bed temperature sensor reading
        sendEvent(name: "temperature", value: data.currentTemperature, unit: "F")
    }
}

// Helper methods
private String getSide() {
    return getDataValue("side") ?: "left"
}

private String getServerIP() {
    return parent.getServerIP()
}

private Integer getServerPort() {
    return parent.getServerPort()
}

private Integer mapTemperatureToLevel(temp) {
    // Linear mapping: 55F = -100, 70F = 0, 110F = 100
    // level = (temp - 70) * (100 / 40) for heating
    // level = (temp - 70) * (100 / 15) for cooling
    def tempF = temp.toFloat()

    if (tempF >= 70) {
        // Heating range: 70F to 110F maps to 0 to 100
        return Math.min(100, ((tempF - 70) * (100 / 40)).toInteger())
    } else {
        // Cooling range: 55F to 70F maps to -100 to 0
        return Math.max(-100, ((tempF - 70) * (100 / 15)).toInteger())
    }
}

private Integer mapLevelToTemperature(level) {
    // Reverse of mapTemperatureToLevel
    def intLevel = level.toInteger()

    if (intLevel >= 0) {
        // Heating: level 0-100 maps to 70-110F
        return 70 + (intLevel * 40 / 100)
    } else {
        // Cooling: level -100 to 0 maps to 55-70F
        return 70 + (intLevel * 15 / 100)
    }
}

private void updateThermostatFromLevel(level) {
    def intLevel = level.toInteger()
    def tempF = mapLevelToTemperature(intLevel)

    if (intLevel >= 0) {
        sendEvent(name: "heatingSetpoint", value: tempF, unit: "F")
        sendEvent(name: "thermostatSetpoint", value: tempF, unit: "F")
        sendEvent(name: "thermostatOperatingState", value: intLevel > 0 ? "heating" : "idle")
    } else {
        sendEvent(name: "coolingSetpoint", value: tempF, unit: "F")
        sendEvent(name: "thermostatSetpoint", value: tempF, unit: "F")
        sendEvent(name: "thermostatOperatingState", value: "cooling")
    }
}

private logDebug(msg) {
    if (parent.getLogEnabled()) {
        log.debug(msg)
    }
}
