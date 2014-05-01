var checked = 0;
function checkReview(ref)
{
  if(checked == 1) {
    return true;
  }
  var x=document.forms["review"]["answer"].value;
  var input=document.getElementsByName("answer")[0];
  var review=document.getElementsByName("checked")[0];
  // show button to show more info!
  var showinfo=document.getElementsByName("showinfobutton")[0];
  showinfo.style.display = "block";
  if (x==ref) {
    input.style.backgroundColor = "limegreen";
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
    return false;
  }
}
function showInfo()
{
  var info=document.getElementsByName("info")[0];
  info.style.display = "block";
}