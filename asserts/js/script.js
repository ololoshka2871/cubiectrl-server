var data_temp, data_spin
var plot_temp, plot_spin

var options = {
		lines: {
			show: true
		},
		points: {
			show: false
		},
		xaxis: {
			show: false
		},
		yaxis: { 
			ticks: 5,
			show: true
		},
		series: {
			shadowSize: 0	// Drawing is faster without shadows
		},
	};

var graphLen = 50;

function update_mesureValues() {	
	$.ajax({
		url: "/data.api",
		dataType: "json",
		method : "POST",
		cache: false,
		data : { "req" : "mesurment", "type" : "temp" },
		success : function(v_data) {

			// create data set
			if 	(typeof data_temp === 'undefined') {	 
				data_temp = new Array();
				for (var serie in v_data) {
					data_temp.push({ label : serie, data : [] })
				}
			} 
			
			plotData(data_temp, v_data, plot_temp, options);
		}
	});
		
	$.ajax({
		url: "/data.api",
		dataType: "json",
		method : "POST",
		cache: false,
		data : { "req" : "mesurment", "type" : "spin" },
		success : function(v_data) {

			// create data set
			if 	(typeof data_spin === 'undefined') {	 
				data_spin = new Array();
				for (var serie in v_data) {
					data_spin.push({ label : serie, data : [] })
				}
			} 
			
			plotData(data_spin, v_data, plot_spin, options);
		}
	});
}

function rereadSettings() {
	$.ajax({
		url: "/data.api",
		dataType: "json",
		method : "POST",
		cache: false,
		data : {"req" : "getSettings"},
		success : function(v_data) {
			if (v_data["OK"] != true)
				alert("Error: " + v_data)
			else {
				for (key in v_data) {
					if (key == "OK") 
						continue
					$("[name='" + key +"']").val(v_data[key])
				}
			}
		}
	});
}

function sendSettings(ev) {
	$.ajax({
		url: "/data.api",
		dataType: "text",
		method : "POST",
		cache: false,
		data : "req=setSettings&" + $("#settingsForm").serialize(),
		success : function(v_data) {
			if (v_data != "OK")
				alert("Error: " + v_data)
			else
				rereadSettings()
		}
	});
}

function resetSettings() {
	$.ajax({
		url: "/data.api",
		dataType: "text",
		method : "POST",
		cache: false,
		data : {"req" : "resetSettings" },
		success : function(v_data) {
			if (v_data != "OK")
				alert("Error: " + v_data)
			else
				rereadSettings()
		}
	});
}
	
function plotData(data, v_data, plot, options) {
	// push new value
	var i = 0;
	for (var item in v_data) {
		var serie = v_data[item];
		data[i].data.push( serie["Error"] ? [null, null] :
				[ Date.parse(serie["Timestamp"]), serie["Value"] ]
		);
		i++;
	}

	// shrink last value
	if (data[0].data.length > graphLen) {
		for (i in data) {
			data[i].data.shift()
		}
	}
	
	// plot
	plot.setData(data);
	plot.setupGrid();
	plot.draw();
}

function playCtrl() {
	var req = {
			req : "DisplayCtrl",
	};
	
	switch ($(this).attr('name')) {
		case "Play_small":
			req.Display = "small";
			ctrl = "play";
			break
		case "Stop_small":
			req.Display = "small";
			ctrl = "stop";
			break
		case "Play_big":
			req.Display = "big";
			ctrl = "play";
			break
		case "Stop_big":
			req.Display = "big";
			ctrl = "stop";
			break
	}
	
	$.ajax({
		url: "/data.api",
		dataType: "text",
		method : "POST",
		cache: false,
		data : req,
		success : function(v_data) {
			if (v_data != "OK")
				alert("Error: " + v_data)
			else
				rereadSettings()
		}
	});
}

jQuery(document).ready(function($){
	// create plots
	plot_temp = $.plot("#plot_temp", [[[]]], options);
	plot_spin = $.plot("#plot_spin", [[[]]], options);
	
	$("#submitBtn").click(sendSettings);
	$("#resetBtn").click(resetSettings);
	
	$("[name*='Play'], [name*='Stop']").click(playCtrl)

	setInterval(update_mesureValues, 500)
});
