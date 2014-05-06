var checked = 0;
function checkReview(ref)
{
  if(checked == 1) {
    return true;
  }
  var x=document.forms["review"]["answer"].value;
  if(x.length === 0) {
    return false;
  }
  var input=document.getElementsByName("answer")[0];
  var review=document.getElementsByName("checked")[0];
  // show button to show more info!
  var showinfo=document.getElementsByName("showinfobutton")[0];
  var bitoff=document.getElementsByName("bitwrong")[0];
  var dist=getBestDistance(x, ref) // dist in %
  if (dist <= 33) { // difference can be 33% :o
    if(dist == 0) {
      input.style.backgroundColor = "limegreen";
      showinfo.style.display = "block";
    } else {
      input.style.backgroundColor = "limegreen";
      bitoff.style.display = "block";
    }
    input.readOnly = true;
    checked = 1;    
    review.value = "true";
    return false;
  } else 
  {
    input.style.backgroundColor = "crimson";
    input.readOnly = true;
    checked = 1;
    review.value = "false";
    showinfo.style.display = "block";
    return false;
  }
};

function getBestDistance(a, ref) {
  var dist = 100;
  for(var b in ref) 
  {
    var tmp = getEditDistance(a.toLowerCase(), ref[b].Imi.toLowerCase());
    dist = Math.min(dist, 100*tmp/ref[b].Imi.length);
  }  
  return dist;
}
// Compute the edit distance between the two given strings
function getEditDistance(a, b) {
  if(a.length === 0) return b.length; 
  if(b.length === 0) return a.length; 
 
  var matrix = [];
 
  // increment along the first column of each row
  var i;
  for(i = 0; i <= b.length; i++){
    matrix[i] = [i];
  }
 
  // increment each column in the first row
  var j;
  for(j = 0; j <= a.length; j++){
    matrix[0][j] = j;
  }
 
  // Fill in the rest of the matrix
  for(i = 1; i <= b.length; i++){
    for(j = 1; j <= a.length; j++){
      if(b.charAt(i-1) == a.charAt(j-1)){
        matrix[i][j] = matrix[i-1][j-1];
      } else {
        matrix[i][j] = Math.min(matrix[i-1][j-1] + 1, // substitution
                                Math.min(matrix[i][j-1] + 1, // insertion
                                         matrix[i-1][j] + 1)); // deletion
      }
    }
  }
 
  return matrix[b.length][a.length];
};
function printTimeLeft()
{
  return "time left: "
}
function showInfo()
{
  var info=document.getElementsByName("info")[0];
  info.style.display = "block";
  document.getElementsByName("answer")[0].focus();
}
function statsChart(stats) {
  var chart = document.getElementById('levelPiechart');
  if(chart.style.display == "block" ) {
    chart.style.display = "none";
  } else {
    chart.style.display = "block";
    var size = Math.min(450, document.getElementById('levelchart').offsetWidth);
    chart.setAttribute('width', size);
    chart.setAttribute('height', size);
    var pieData = [
    {
      value: stats.Beginner,
      color : "#F7464A"
    },
    {
      value : stats.Elementary,
      color : "#F38630"
    },
    {
      value : stats.Intermediate,
      color : "#F0E68C"
    },
    {
      value : stats.Master,
      color : "#48D1CC"
    },
    {
      value : stats.Known,
      color : "#69D2E7"
    }];
    var myPie = new Chart(chart.getContext("2d")).Pie(pieData);
  }
}