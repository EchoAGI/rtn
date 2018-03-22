"use strict";


define(['api', 'connector'], function(api, connector) {
    var channelClient = function(version, url, config) {
        var me = this
        me.version = version || "1.0"
        me.url = url || "wss://127.0.0.1:8443"
        me.config = config

        me.connector = new Connector()

        me.api = new Api("1.0", connector)


        me.connector.e.on("open error close", function(event) {
			switch (event.type) {
				case "open":
					flags.connected = true
					flags.autoreconnectDelay = 0
					me.api.updateStatus(true)
					break;
				case "error":
					if (config.connected) {
						reconnect()
					} else {
					}
					break;
				case "close":
					reconnect()
					break;
			}
		});



        this.connector.connect(url)
    }


    var reconnect = function() {
        if (appData.flags.connected && appData.flags.autoreconnect) {
            if (appData.flags.resurrect === null) {
                // Store data at the resurrection shrine.
                appData.flags.resurrect = {
                    status: $scope.getStatus(),
                    id: $scope.id
                }
                console.log("Stored data at the resurrection shrine", appData.flags.resurrect);
            }
            if (!appData.flags.reconnecting) {
                var delay = appData.flags.autoreconnectDelay;
                if (delay < 10000) {
                    appData.flags.autoreconnectDelay += 500;
                }
                appData.flags.reconnecting = true;
                _.delay(function() {
                    if (appData.flags.autoreconnect) {
                        console.log("Requesting to reconnect ...");
                        mediaStream.reconnect();
                    }
                    appData.flags.reconnecting = false;
                }, delay);
                $scope.setStatus("reconnecting");
            } else {
                console.warn("Already reconnecting ...");
            }
        } else {
            $scope.setStatus("closed");
        }
    };
})