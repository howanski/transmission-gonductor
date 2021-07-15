(function (window, document, $) {

    let gonductorSettings = {};

    let formInputsSetup = {
        // format -"settingsKey": "selector"
        "transmissionHost": "#inputTransmissionHost",
        "transmissionUser": "#inputTransmissionUser",
        "transmissionPassword": "#inputTransmissionPassword",
    };

    let formCheckboxesSetup = {
        // format -"settingsKey": "selector"
        "transmissionManagePrioritiesAlphabetically": "checkTransmissionAlphabeticalPriority",
        "gonductorDebugToTerminal": "gonductorSpamTerminal"
    }

    $(document).ready(function () {
        console.log("--------------------------------------");
        console.log("jQuery initialized, running GOnductor.");
        console.log("--------------------------------------");
        readSettings();
        runInLoop(statsReader, 5);
        $('.saveBtn').on('click', function (e) {
            e.preventDefault();
            saveSettings();
        });
    });

    let statsReader = function () {
        let statisticsMap = {
            "connectionStatus": "#statConnectionStatus",
            "lastPing": "#statLastPing"
        };
        $.ajax({
            type: "GET",
            url: '/gonductor-stats',
            success: function (response) {
                $.each(response, function (index, value) {
                    let statBox = $(statisticsMap[index]);
                    statBox.html(value);
                });
            },
            error: function (response){
                $.each(statisticsMap, function (index, selector) {
                    let statBox = $(selector);
                    statBox.html('GOnductor server GOne');
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
                updateFormDataFromGlobalSettings();
                console.log("-----------------------------");
                console.log("GOnductor settings (re)loaded");
                console.log("-----------------------------");
            }
        });
    }

    let saveSettings = function () {
        readFormDataToGlobalSettings();
        $.ajax({
            type: "POST",
            url: '/settings',
            data: JSON.stringify(gonductorSettings),
            success: function (response) {
                console.log("------------------------");
                console.log("GOnductor settings saved");
                console.log("------------------------");
                readSettings();
            },
            error: function (response) {
                alert('Save error');
            }
        });
    }

    let readFormDataToGlobalSettings = function () {
        $.each(formInputsSetup, function (key, selector) {
            let jqInput = $(selector);
            gonductorSettings[key] = jqInput.val();
        });
        $.each(formCheckboxesSetup, function (key, selectorId) {
            let jqInput = $("#" + selectorId);
            if (jqInput.is(":checked")) {
                gonductorSettings[key] = "yes";
              } else {
                gonductorSettings[key] = "no";
              }
        });
    }

    let updateFormDataFromGlobalSettings = function () {
        $.each(formInputsSetup, function (key, selector) {
            let jqInput = $(selector);
            jqInput.val(gonductorSettings[key]);
        });
        $.each(formCheckboxesSetup, function (key, selectorId) {
            // jQuery sucks in checkboxes, attr, nor prop are working
            let cBox = document.getElementById(selectorId);
            if (gonductorSettings[key] == "yes"){
                cBox.checked = true;
            } else {
                cBox.checked = false;
            }
        });
    }

    let runInLoop = function (func, loopSeconds) {
        func();
        setInterval(func, loopSeconds * 1000);
    }

})(window, document, jQuery);