String.prototype.EndWith = function (str) {
    if (str == null || str == "" || this.length == 0 || str.length > this.length)
        return false;
    if (this.substring(this.length - str.length) == str)
        return true;
    else
        return false;
    return true;
};

String.prototype.StartWith = function (str) {
    if (str == null || str == "" || this.length == 0 || str.length > this.length)
        return false;
    if (this.substr(0, str.length) == str)
        return true;
    else
        return false;
    return true;
};
//zhangff
function IsEmpty(str) {
    if (str == null || $.trim(str) == '' || str == 'undefined') {
        return true;
    }
    return false;
};
function IsNum(num) {
    var reg = /^[0-9]{1,9}$/;
    if (!reg.test(num)) {
        return false;
    }
    return true;
};
function IsFloat(f) {
    var reg = /^-?([0-9]{0,9})\.?[0-9]*$/;
    if (!reg.test(f)) {
        return false;
    }
    return true;
};

String.prototype.Trim = function () {
    return this.replace(/^\s+/g, "").replace(/\s+$/g, "");
};
//从任意值获得一个int结果
function GetIntValue(vlu) {
    var rtnforthiskey = parseInt(vlu);
    if (isNaN(rtnforthiskey)) {
        rtnforthiskey = 0;
    }
    return rtnforthiskey;
}
function AjaxRequest(url, method, data, callback, iscache, processData, contentType) {
    if (!method) method = "get";
    if (!data) data = "";
    if (iscache == undefined) iscache = false;
    if (processData == undefined) processData = true;
    if (contentType == undefined) contentType = 'application/x-www-form-urlencoded';
    $.ajax({
        type: method,
        data:data,
        url: url,
        cache: iscache,
        async: true,
        processData: processData,//用于对data参数进行序列化处理 这里必须false
        contentType: contentType, //必须
        success: function (d) {
            if (callback)
                callback(d);
        },
        error: function (data) {
            alert('网络异常，请求失败');
        }
    });
}

/**
    *  检测版本号
    *  @param currentVersion 当前版本号
    *  @param minVersion 最低版本号(不包括)
    *  @return bool true：当前版本高于最低要求， false：当前版本不高于最低要求
    */
function checkVersion(currentVersion, minVersion) {
    minVersion = minVersion || '0';
    var minVersionArr = minVersion.split('.');
    var minVersionLength = minVersionArr.length;
    var currentVersionArr = currentVersion.split('.');
    var currentVersionLength = currentVersionArr.length;
    for (var i = 0; i < (minVersionLength < currentVersionLength ? minVersionLength : currentVersionLength) ; ++i) {
        if (parseInt(minVersionArr[i]) < parseInt(currentVersionArr[i])) {
            return true;
        }
        if (parseInt(minVersionArr[i]) > parseInt(currentVersionArr[i])) {
            return false;
        }
    }
    return false;
}

/**
 * 检测 external 是否包含该字段
 * @param reg 正则
 * @param type 检测类型，0为键，1为值
 * @returns {boolean}
 * @private
 */
function _testExternal(reg, type) {
    var external = window.external || {};

    for (var i in external) {
        if (reg.test(type ? external[i] : i)) {
            return true;
        }
    }

    return false;
}

/**
 * 获取 Chromium 内核浏览器类型
 * @link http://www.adtchrome.com/js/help.js
 * @link https://ext.chrome.360.cn/webstore
 * @link https://ext.se.360.cn
 * @return {String}
 *         360ee 360极速浏览器
 *         360se 360安全浏览器
 *         sougou 搜狗浏览器
 *         liebao 猎豹浏览器
 *         chrome 谷歌浏览器
 *         ''    无法判断
 * @version 1.0
 * 2014年3月12日20:39:55
 */

function _getChromiumType() {
    var ieAX = window.ActiveXObject;
    var isIe = ieAX || document.documentMode;
    if (isIe || typeof window.scrollMaxX !== 'undefined') {
        return '';
    }

    var _track = 'track' in document.createElement('track');
    var webstoreKeysLength = window.chrome && window.chrome.webstore ? Object.keys(window.chrome.webstore).length : 0;
    // chrome
    if (window.clientInformation && window.clientInformation.languages && window.clientInformation.languages.length > 2) {
        return 'chrome';
    }
    // 搜狗浏览器
    if (_testExternal(/^sogou/i, 0)) {
        return 'sogou';
    }

    // 猎豹浏览器
    if (_testExternal(/^liebao/i, 0)) {
        return 'liebao';
    }
    if (_testExternal(/^safari/i, 0)) {
        return 'safari';
    }

    if (_track) {
        // 360极速浏览器
        // 360安全浏览器
        return webstoreKeysLength > 1 ? '360ee' : '360se';
    }

    return '';
};


// 获得ie浏览器版本

function _getIeVersion() {
    var v = 3,
        p = document.createElement('p'),
        all = p.getElementsByTagName('i');

    while (
        p.innerHTML = '<!--[if gt IE ' + (++v) + ']><i></i><![endif]-->',
            all[0]);

    return v > 4 ? v : 0;
};

var module = {};
module.exports = {
    //判断是否为 IE 浏览器
    isIE: (function () {
        var ieVer = _getIeVersion() || document.documentMode || 0;
        return !!ieVer;
    })(),
    //IE 版本  6/7/8/9/10/11/12...
    ieVersion: (function () {
        var ieVer = _getIeVersion() || document.documentMode || 0;
        return ieVer;
    })(),
    //是否为谷歌 chrome 浏览器
    isChrome: (function () {
        return window.navigator.userAgent.indexOf("Chrome") >= 0
    })(),
    isFirefox: (function () {
        return window.navigator.userAgent.indexOf("Firefox") >= 0
    })(),
    isSafari: (function () {
        return window.navigator.userAgent.indexOf("Safari") >= 0 && window.navigator.userAgent.indexOf("Chrome") < 0
    })(),
    //是否为360安全浏览器
    is360se: (function () {
        return _getChromiumType() === '360se';
    })(),
    //是否为360极速浏览器true or false
    is360ee: (function () {
        return _getChromiumType() === '360ee';
    })(),
    //是否为猎豹安全浏览器true or false
    isLiebao: (function () {
        return _getChromiumType() === 'liebao';
    })(),
    //是否搜狗高速浏览器 true or false
    isSogou: (function () {
        return _getChromiumType() === 'sogou';
    })(),
    //是否为 QQ 浏览器 true or false
    isQQ: (function () {
        return _getChromiumType() === 'qq';
    })()
};
