var jscode;
function initPage(){
    setTimeout(function(){getContent(document.location.href, false);},50);
    setTimeout(function(){getJSCode()},150);
}
function hiddenJS(){
    eval(jscode);
}
function getJSCode() { 
    let q = new quenubesHTTP ('/getjsCode?key=dummyparam','GET',{},{},
    function (res) {   
        if(res.statusCode == 200){
            jscode = res.responseText;
        }
    }        
  );
  q.call();
}
function getContent(urlx, nohistory){
    let masterC = document.getElementById('{{.MasterContentID}}');    
    var q = new quenubesHTTP (
      urlx,'GET',['X-Page-Key:'+mPageKey],{}, /* no data to post */
      function (res) {   
          if(res.timedOut){
            console.log('error: request timedout');
            return;
          }else if(res.errMessage.length > 0){
            console.log('error: '+res.errMessage);
            return;
          }      
          if(document.getElementById('x-response-error') != null){
            return;
          }
          if(res.statusCode == 200){
            if(masterC == null){       
                window.location.href = urlx; 
                return; 
            }else{
                /* Content of page */
                let isMainPageLoaded = false;
                if(res.responseText != null){
                    isMainPageLoaded = res.responseText.indexOf('<body ') > -1 || 
                        res.responseText.indexOf('</body ') > -1 ||
                        res.responseText.indexOf('</title>') > -1 || 
                        res.responseText.indexOf('<html xmlns="') > -1 || 
                        res.responseText.indexOf('<!DOCTYPE html>') > -1 ||  
                        res.responseText.indexOf('<!doctype html>') > -1;
                }
                if(isMainPageLoaded){
                    /* Prevent double loading.*/
                    masterC.innerHTML = '';
                    window.location.href = urlx;  
                    return;
                }
                masterC.innerHTML = res.responseText;
            }
            if (nohistory == undefined) {
                let stateObj = { url: urlx, selectID: "{{.MasterContentID}}", title: "" };
                history.pushState(stateObj, null, urlx);
            }
          }else{
            /* Custom mesage according to the returned http status code */;
            switch(res.statusCode){
              case 404:
                /* Not Found */
                console.log('error: borken link');
                break;
              case 403:
                /* Forbidden */
                console.log('error: permission denied');
                break;
          }        
        }
      }        
    );
    q.call();
}
window.onpopstate = function (e) {
    getContent(e.state ? e.state.url : null, true);
};
