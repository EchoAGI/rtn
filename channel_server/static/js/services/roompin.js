"use strict";
define([
], function() {

	return ["$window", "$q", "alertify", "translation", function($window, $q, alertify, translation) {

		var pinCache = {};
		var roompin = {
			get: function(roomName) {
				var cachedPIN = pinCache[roomName];
				return cachedPIN ? cachedPIN : null;
			},
			clear: function(roomName) {
				delete pinCache[roomName];
				console.log("Cleared PIN for", roomName);
			},
			update: function(roomName, pin) {
				if (pin) {
					pinCache[roomName] = pin;
					alertify.dialog.alert(translation._("PIN for room %s is now '%s'.", roomName, pin));
				} else {
					roompin.clear(roomName);
					alertify.dialog.alert(translation._("PIN lock has been removed from room %s.", roomName));
				}
			},
			requestInteractively: function(roomName) {
				var deferred = $q.defer();
				alertify.dialog.prompt(translation._("Enter the PIN for room %s", roomName), function(pin) {
					if (pin) {
						pinCache[roomName] = pin;
						deferred.resolve();
					} else {
						deferred.reject();
					}
				}, function() {
					deferred.reject();
				});
				return deferred.promise;
			}
		};

		return roompin;

	}];
});
