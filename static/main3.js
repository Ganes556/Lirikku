import{h as t}from"./assets/htmx.min-BDSWTZMZ.js";(function(){var i,d="hx-target-";function u(e,r){return e.substring(0,r.length)===r}function l(e,r){if(!e||!r)return null;var s=r.toString(),a=[s,s.substr(0,2)+"*",s.substr(0,2)+"x",s.substr(0,1)+"*",s.substr(0,1)+"x",s.substr(0,1)+"**",s.substr(0,1)+"xx","*","x","***","xxx"];(u(s,"4")||u(s,"5"))&&a.push("error");for(var o=0;o<a.length;o++){var g=d+a[o],f=i.getClosestAttributeValue(e,g);if(f)return f==="this"?i.findThisElement(e,g):i.querySelectorExt(e,f)}return null}function n(e){e.detail.isError?t.config.responseTargetUnsetsError&&(e.detail.isError=!1):t.config.responseTargetSetsError&&(e.detail.isError=!0)}t.defineExtension("response-targets",{init:function(e){i=e,t.config.responseTargetUnsetsError===void 0&&(t.config.responseTargetUnsetsError=!0),t.config.responseTargetSetsError===void 0&&(t.config.responseTargetSetsError=!1),t.config.responseTargetPrefersExisting===void 0&&(t.config.responseTargetPrefersExisting=!1),t.config.responseTargetPrefersRetargetHeader===void 0&&(t.config.responseTargetPrefersRetargetHeader=!0)},onEvent:function(e,r){if(e==="htmx:beforeSwap"&&r.detail.xhr&&r.detail.xhr.status!==200){if(r.detail.target&&(t.config.responseTargetPrefersExisting||t.config.responseTargetPrefersRetargetHeader&&r.detail.xhr.getAllResponseHeaders().match(/HX-Retarget:/i)))return r.detail.shouldSwap=!0,n(r),!0;if(!r.detail.requestConfig)return!0;var s=l(r.detail.requestConfig.elt,r.detail.xhr.status);return s&&(n(r),r.detail.shouldSwap=!0,r.detail.target=s),!0}}})})();
