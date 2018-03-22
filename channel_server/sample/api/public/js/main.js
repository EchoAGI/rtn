"use strict";


require.config({
	waitSeconds: 300,
	paths: {
		// Major libraries
		"jquery": 'libs/jquery/jquery.min',
		"underscore": 'libs/lodash.min', // alternative to underscore
		'ua-parser': 'libs/ua-parser'
	},
	shim: {
		'underscore': {
			exports: '_'
		}
	}
});


require.onError = (function() {
	return function(err) {
		if (err.requireType === "timeout" || err.requireType === "scripterror") {
			console.error("Error while loading " + err.requireType, err.requireModules);
		} else {
			throw err;
		}
	};
}());



require(['api', 'connector'], function(Api, Connector) {
    var context = {
        Host: "127.0.0.1",
        Port: 8443,
    }

    var connector = new Connector()
    var api = new Api("1.0", connector)

    var url = "wss" + "://" + context.Host + ":" + context.Port + "/ws";
    console.log("url: " + url)
    connector.connect(url)


    api.e.on("received.self", function(event, data) {
        alert("received.self")
    
        var roomName = "haha"
        var pin = "123"
        api.sendJoinRoom(roomName, pin, function(room) {
            console.log("room: ", room)
        }, function(error) {
            console.error(error)
        });
    })

    api.e.on("received.chat", function(event, id, from, data, p2p) {
        console.log("event: ", event)
        console.log("id: ", id)
        console.log("from: ", from)
        console.log("data: ", data)
        console.log("p2p: ", p2p)
    })

    api.e.on("received.users", function(event, data) {
        _.each(data, function(p) {
            console.log('user: ', p)
        });
    });
});