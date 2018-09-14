(function ($) {
    $.fn.extend({
        initForm: function (options) {
            //默认参数  
            var defaults = {
                jsonValue: "",
                isDebug: false //是否需要调试，这个用于开发阶段，发布阶段请将设置为false，默认为false,true将会把name value打印出来  
            }
            //设置参数  
            var setting = $.extend({}, defaults, options);
            var form = this;

            //如果传入的json对象为空，则不做任何操作  
            if (!$.isEmptyObject(options)) {
                var debugInfo = "";
                $.each(options,
                    function (key, value) {
                        //是否开启调试，开启将会把name value打印出来  
                        if (setting.isDebug) {
                            alert("name:" + key + "; value:" + value);
                            debugInfo += "name:" + key + "; value:" + value + " || ";
                        }
                        var formField = form.find("[name='" + key + "']");
                        if ($.type(formField[0]) === "undefined") {
                            if (setting.isDebug) {
                                alert("can not find name:[" + key + "] in form!!!"); //没找到指定name的表单  
                            }
                        } else {
                            var fieldTagName = formField[0].tagName.toLowerCase();
                            if (fieldTagName == "input") {
                                if (formField.attr("type") == "radio") {
                                    $("input:radio[name='" + key + "'][value='" + value + "']")
                                        .attr("checked", "checked");
                                } else {
                                    formField.val(value);
                                }
                            } else if (fieldTagName == "select") {
                                //do something special  
                                formField.val(value);
                            } else if (fieldTagName == "textarea") {
                                //do something special  
                                formField.val(value);
                            }
                            else if (fieldTagName == "label") {
                                //do something special  
                                formField.html(value);
                            } else {
                                formField.val(value);
                            }
                        }
                    });
                if (setting.isDebug) {
                    alert(debugInfo);
                }
            }
            return form; //返回对象，提供链式操作  
        },
        disableForm: function (formobj, isDisabled) {

            for (var i = 0; i < formobj.length; i++) {
                var element = formobj.elements[i];
                element.disabled = isDisabled;
            }
        }
    });
})(jQuery);

