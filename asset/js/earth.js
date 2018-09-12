function startEarth(data) {
	var show = false;

	function initialize() {
		var options = { zooming:false, sky:true}
		if (window.location.href.substr(0, 3) === 'file:')
			alert("This file must be accessed via http:// or https:// to run properly.");
		var earth = new WE.map('earth_div', options);
		earth.setView([46.473848, 30.711033], 5);
		WE.tileLayer('img/{z}/{x}/{y}.png', {
			tileSize: 256,
			bounds: [[-85, -180], [85, 180]],
			attribution: '',
			scrollWheelZoom:false
		}).addTo(earth);
		

		var before = null;
		
		// requestAnimationFrame(function animate(now) {
		// 	if (!show) return;

		// 	var c = earth.getPosition();
		// 	var elapsed = before? now - before: 0;
		// 	before = now;
		// 	earth.setCenter([c[0], c[1] + 0.01*(elapsed/60)]);
		// 	earth.setTilt(40)
		// 	requestAnimationFrame(animate);
		// });
		earth.setTilt(40)

		for(var i = 0; i < data.length; i++) {
			console.log("Init marker: " + data[i].x + ", " + data[i].y + " " + data[i].name)

			var marker = WE.marker([data[i].x, data[i].y], false).addTo(earth)
			marker.bindPopup('<div id="e-chat'+ data[i].chat_id +'" class="message-home"></div>').openPopup ();
		}
		return earth;	
	}


	setInterval(function() {
			var epos = $("#earth_div").offset().top,
				upos = $(window).scrollTop(),
				eheight = $("#earth_div").height(),
				uheight = $(window).height();

			if (upos >= epos-(uheight*2) && upos <= epos+eheight) {
				if (!show) {
					initialize()
					console.log("earth show")
				}
				show = true;
			} else {

				if(show) {
					$("#earth_div").html("")
					console.log("earth hide")
				}
				show = false;
			}
	}, 1000)
}
