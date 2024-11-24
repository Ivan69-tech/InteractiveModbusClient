
<script src="https://cdn.jsdelivr.net/npm/chart.js"></script>

    // Remplacer ces valeurs avec vos données
    // const timestamps = [
    //     '2024-11-21T00:00:00Z', 
    //     '2024-11-21T01:00:00Z', 
    //     '2024-11-21T02:00:00Z',
    //     '2024-11-21T03:00:00Z'
    // ]; // Liste de timestamps
    //const data = [10, 20, 15, 25]; // Liste de valeurs entières
    console.log("ok")
    function fetchData() {
        fetch('/getdb')
            .then(response => {
                if (response.ok) {
                    return response.json();
                } else {
                    throw new Error('Request failed with status: ' + response.status);
                }
            })
            .then(data => {
                console.log(data);
            })
            .catch(error => {
                console.error('Request failed', error);
            });
    }
    
    
    // const myChart = new Chart(ctx, {
    //     type: 'line',
    //     data: {
    //         labels: timestamps,
    //         datasets: [{
    //             label: 'Sample Data',
    //             data: data,
    //             backgroundColor: 'rgba(75, 192, 192, 0.2)',
    //             borderColor: 'rgba(75, 192, 192, 1)',
    //             borderWidth: 1,
    //             fill: false
    //         }]
    //     },
    //     options: {
    //         scales: {
    //             x: {
    //                 type: 'time',
    //                 time: {
    //                     unit: 'hour',
    //                     tooltipFormat: 'll HH:mm'
    //                 },
    //                 title: {
    //                     display: true,
    //                     text: 'Timestamp'
    //                 }
    //             },
    //             y: {
    //                 beginAtZero: true,
    //                 title: {
    //                     display: true,
    //                     text: 'Value'
    //                 }
    //             }
    //         }
    //     }
    // });