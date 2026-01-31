/**
 * Eight Sleep Pod
 * Parent device driver for Eight Sleep Pod integration
 * Requires eightctl hubitat server running locally
 *
 * Copyright 2024 eightctl
 * Licensed under the Apache License, Version 2.0
 */

metadata {
    definition(name: "Eight Sleep Pod", namespace: "eightctl", author: "eightctl") {
        capability "Refresh"
        capability "Initialize"

        command "createChildDevices"
        command "deleteChildDevices"

        attribute "connectionStatus", "string"
        attribute "lastRefresh", "string"
    }

    preferences {
        input name: "serverIP", type: "string", title: "Server IP", description: "IP address of eightctl hubitat server", defaultValue: "192.168.1.x", required: true
        input name: "serverPort", type: "number", title: "Server Port", description: "Port of eightctl hubitat server", defaultValue: 8080, required: true
        input name: "refreshInterval", type: "number", title: "Refresh Interval (seconds)", description: "How often to poll for status updates", defaultValue: 60, required: true
        input name: "logEnable", type: "bool", title: "Enable debug logging", defaultValue: false
    }
}

def installed() {
    logDebug("Eight Sleep Pod installed")
    initialize()
}

def updated() {
    logDebug("Eight Sleep Pod updated with settings: serverIP=${serverIP}, serverPort=${serverPort}, refreshInterval=${refreshInterval}")
    unschedule()
    initialize()
}

def initialize() {
    logDebug("Initializing Eight Sleep Pod")
    sendEvent(name: "connectionStatus", value: "initializing")

    if (!serverIP || serverIP == "192.168.1.x") {
        log.warn("Server IP not configured - please set the eightctl hubitat server IP address")
        sendEvent(name: "connectionStatus", value: "not configured")
        return
    }

    // Schedule periodic refresh
    if (refreshInterval && refreshInterval > 0) {
        def cronExpression = "0/${Math.max(refreshInterval, 10)} * * * * ?"
        schedule(cronExpression, refresh)
        logDebug("Scheduled refresh every ${refreshInterval} seconds")
    }

    // Initial refresh
    runIn(2, refresh)
}

def refresh() {
    logDebug("Refreshing Eight Sleep Pod status")

    def params = [
        uri: "http://${serverIP}:${serverPort}",
        path: "/status",
        contentType: "application/json",
        timeout: 10
    ]

    try {
        httpGet(params) { response ->
            if (response.status == 200) {
                def data = response.data
                logDebug("Received status: ${data}")

                sendEvent(name: "connectionStatus", value: "connected")
                sendEvent(name: "lastRefresh", value: new Date().format("yyyy-MM-dd HH:mm:ss"))

                // Update child devices
                updateChildDevice("left", data.left)
                updateChildDevice("right", data.right)
            } else {
                log.error("Unexpected response status: ${response.status}")
                sendEvent(name: "connectionStatus", value: "error")
            }
        }
    } catch (groovyx.net.http.HttpResponseException e) {
        log.error("HTTP error refreshing status: ${e.message}")
        sendEvent(name: "connectionStatus", value: "error")
    } catch (java.net.ConnectException e) {
        log.error("Connection refused - is eightctl hubitat server running at ${serverIP}:${serverPort}?")
        sendEvent(name: "connectionStatus", value: "connection refused")
    } catch (Exception e) {
        log.error("Error refreshing status: ${e.message}")
        sendEvent(name: "connectionStatus", value: "error")
    }
}

def createChildDevices() {
    logDebug("Creating child devices for left and right sides")

    ["left", "right"].each { side ->
        def childDni = "${device.deviceNetworkId}-${side}"
        def existingChild = getChildDevice(childDni)

        if (!existingChild) {
            logDebug("Creating child device for ${side} side")
            try {
                def child = addChildDevice(
                    "eightctl",
                    "Eight Sleep Side",
                    childDni,
                    [
                        name: "Eight Sleep ${side.capitalize()}",
                        label: "${device.label ?: device.name} - ${side.capitalize()}",
                        isComponent: true
                    ]
                )
                child.updateDataValue("side", side)
                log.info("Created child device: ${child.label}")
            } catch (Exception e) {
                log.error("Error creating child device for ${side}: ${e.message}")
            }
        } else {
            logDebug("Child device for ${side} side already exists")
        }
    }

    // Refresh to populate child device states
    runIn(2, refresh)
}

def deleteChildDevices() {
    logDebug("Deleting all child devices")
    getChildDevices().each { child ->
        log.info("Deleting child device: ${child.label}")
        deleteChildDevice(child.deviceNetworkId)
    }
}

def updateChildDevice(String side, Map sideData) {
    if (!sideData) {
        logDebug("No data for ${side} side")
        return
    }

    def childDni = "${device.deviceNetworkId}-${side}"
    def child = getChildDevice(childDni)

    if (!child) {
        logDebug("Child device for ${side} side not found - run createChildDevices first")
        return
    }

    logDebug("Updating ${side} side with data: ${sideData}")
    child.updateStatus(sideData)
}

// Called by child devices to get server settings
String getServerIP() {
    return serverIP
}

Integer getServerPort() {
    return serverPort
}

Boolean getLogEnabled() {
    return logEnable
}

private logDebug(msg) {
    if (logEnable) {
        log.debug(msg)
    }
}
