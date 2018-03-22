"use strict";

define([
	'mediastream/connector'
], function(Connector) {
	return [function() {
		return new Connector();
	}];
});
