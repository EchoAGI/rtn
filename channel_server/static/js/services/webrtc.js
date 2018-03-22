"use strict";
define([
	'mediastream/webrtc'
], function(WebRTC) {
	return ["api", function(api) {
		return new WebRTC(api);
	}];
});
