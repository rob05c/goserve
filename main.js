function getXmlHttpRequest() {
    if(window.XMLHttpRequest) {
	return new XMLHttpRequest();
    }
    return new ActiveXObject("Microsoft.XMLHTTP");
}

function loadPage(e, uri, contentId) {

    var sender = (e && e.target) || (window.event && window.event.srcElement);
    sender = sender.parentNode;
    var siblings = sender.parentNode.children;
    for(var i = 0; i != siblings.length; i++) {
	siblings[i].className = "";
    }
    sender.className = "active";

    var xmlhttp = getXmlHttpRequest();
    xmlhttp.onreadystatechange = function() {
	if(xmlhttp.readyState == 4 && xmlhttp.status == 200) {
	    document.getElementById(contentId).innerHTML = xmlhttp.responseText;
	}
    }
    xmlhttp.open("GET", uri, true);
    xmlhttp.send();
}
