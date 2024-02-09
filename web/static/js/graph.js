document.addEventListener('DOMContentLoaded', DefaultConfig);
document.getElementById('DefaultConfig').addEventListener('click',DefaultConfig);
for (let i = 0; i < 12; i++) {
    const display = document.getElementById(`sliderValue${i}`);
    const slider = document.getElementById(`slider${i}`).oninput = function() {
        display.innerHTML = this.value;
        GetUpdatedPositions();
    };
    display.innerHTML = slider.value;
}

function SendRequest(url, configure, func) {
    fetch(url, func)
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json(); // Assuming the server responds with JSON
        })
        .then(data => {
            UpdateGraph(data["quadrubot"])
            if (configure) {
                for (let i = 0; i < data["angles"].length; i++) {
                    for (let j = 0; j < 3; j++){
                        const slider = document.getElementById(`slider${i * 3 + j}`)
                        slider.value = data["angles"][i][j]
                        document.getElementById(`sliderValue${i * 3 + j}`).innerHTML = slider.value;
                    }
                }
            }
        })
        .catch(error => {
            console.error('Fetch error:', error);
        });
}

function DefaultConfig() {
    SendRequest('DefaultConfig', true, {
        method: 'GET',
    })
}

function GetUpdatedPositions() {
    const sliderValues = Array.from(document.querySelectorAll('.slider-input')).map(input => input.value);
    SendRequest('/', false, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({sliderValues: sliderValues})
    })
}

function UpdateGraph(data) {
    let quadrupletOptions = {
        "animation": true,
        "grid3D": {},
        "legend": {"show": true, "type": ""},
        "series": [],
        "title": {},
        "tooltip": {"show": true},
        "xAxis3D": {"min": -500, "max": 500},
        "yAxis3D": {"min": -500, "max": 500},
        "zAxis3D": {"min": -500, "max": 500}
    };
    for (let i = 0; i < data.length; i++) {
        quadrupletOptions.series.push({
            "name": data[i].Name,
            "type": "line3D",
            "smooth": true,
            "connectNulls": false,
            "showSymbol": true,
            "waveAnimation": false,
            "coordinateSystem": "cartesian3D",
            "renderLabelForZeroData": false,
            "selectedMode": false,
            "animation": true,
            "data": data[i].Data
        })
    }
    let quadrupletChart = echarts.init(document.getElementById('quadruplet'), "westeros");
    quadrupletChart.setOption(quadrupletOptions);
}
