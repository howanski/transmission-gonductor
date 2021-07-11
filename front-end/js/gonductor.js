(function (window, document, $) {
    $(document).ready(function () {
        console.log("------------------------------------");
        console.log("jQuery initialized, GOnductor ready.");
        console.log("------------------------------------");
        runInLoop(statsReader, 5);
        console.log("-------------------------------");
        console.log("GOnductor functions initialized");
        console.log("-------------------------------");
    });

    let statsReader = function () {
        $.ajax({
            type: "GET",
            url: '/gonductor-stats',
            success: function (response) {
                let statisticsMap = {
                    "connectionStatus": "#statConnectionStatus",
                    "lastPing": "#statLastPing"
                };
                $.each(response, function (index, value) {
                    let statBox = $(statisticsMap[index]);
                    statBox.html(value);
                });
            }
        });
    };

    let runInLoop = function (func, loopSeconds) {
        func();
        setInterval(func, loopSeconds * 1000);
    }

})(window, document, jQuery);