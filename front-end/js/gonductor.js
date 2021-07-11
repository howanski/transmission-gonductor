(function (window, document, $) {

    let gonductorSettings = {};

    $(document).ready(function () {
        console.log("--------------------------------------");
        console.log("jQuery initialized, running GOnductor.");
        console.log("--------------------------------------");
        readSettings();
        runInLoop(statsReader, 5);
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

    let readSettings = function () {
        $.ajax({
            type: "GET",
            url: '/settings',
            success: function (response) {
                $.each(response, function (index, settingsObject) {
                    gonductorSettings[settingsObject.ConfigKey] = settingsObject.ConfigValue;
                });
                console.log("-------------------------");
                console.log("GOnductor settings loaded");
                console.log("-------------------------");
            }
        });
    }

    let runInLoop = function (func, loopSeconds) {
        func();
        setInterval(func, loopSeconds * 1000);
    }

})(window, document, jQuery);