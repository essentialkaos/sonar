function sonarStatusReferesher() {
  var points = document.getElementsByClassName("slack-status");

  if (points.length == 0) {
    return;
  }

  for (var i = 0; i < points.length; i++) {
    var point = points[i];
    var pointURI = point.src;

    var rnd = Math.round(Math.random() * 1000000000000);
    var marker = pointURI.indexOf("&rnd=");
    var wasPointRefreshedEarlier = marker !== -1

    if (wasPointRefreshedEarlier) {
      point.src = point.src.substr(0, marker) + "&rnd=" + rnd.toString(26);
    } else {
      point.addEventListener("error", function() {
        this.style.visibility = "hidden";
      });

      point.addEventListener("load", function() {
        this.style.visibility = "visible";
      });
      
      point.src = point.src + "&rnd=" + rnd.toString(26);
    }
  }
}

setInterval(sonarStatusReferesher, 120000);
