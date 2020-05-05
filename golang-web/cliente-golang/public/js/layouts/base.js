const cabeceras= {
    'Content-Type': 'application/json',
    'Accept': 'application/json',
}

function getCookie(cname) {
    var name = cname + "=";
    var decodedCookie = decodeURIComponent(document.cookie);
    var ca = decodedCookie.split(';');
    for(var i = 0; i <ca.length; i++) {
      var c = ca[i];
      while (c.charAt(0) == ' ') {
        c = c.substring(1);
      }
      if (c.indexOf(name) == 0) {
        return c.substring(name.length, c.length);
      }
    }
    return "";
}

function checkCookie(cname) {
    var user = getCookie(cname);
    if (user != "") {
      return true;
    } else {
      return false;
    }
}

function delete_cookie(cname) {
    document.cookie = cname +'=; Path=/; Expires=Thu, 01 Jan 1970 00:00:01 GMT;';
  }

function displayLogoutButton(){
    if(checkCookie("userSession")){
        document.querySelector("#userMenuDropdown").classList.remove("invisible");
        document.querySelector("#loginMenu").classList.add("invisible");
        document.querySelector("#registerMenu").classList.add("invisible");
    }else{
      document.querySelector("#userMenuDropdown").classList.add("invisible");
      document.querySelector("#loginMenu").classList.remove("invisible");
      document.querySelector("#registerMenu").classList.remove("invisible");
    }
}

function logout(){
    const url= `/logout`;
    const payload= {};
    const request = {
        method: 'POST', 
        headers: cabeceras,
        body: JSON.stringify(payload),
    };
    fetch(url,request)
    .then( console.log("SESION CERRADA"), delete_cookie("userSession"), window.location.href="/home")
    .catch(response => alert(response))
}


//PARA EL BREADCRUMB
function deleteBreadcrumb(){
  var breadcrump = document.querySelector("#breadcrumpMenu");
  while (breadcrump.firstChild) {
    breadcrump.removeChild(breadcrump.lastChild);
  }
}

function addLinkBreadcrumb(nombre, url){
  var breadcrump = document.querySelector("#breadcrumpMenu");
  let li = document.createElement('li');
  li.classList.add("breadcrumb-item");
  let a = document.createElement('a');
  a.textContent = nombre;
  a.setAttribute('href', url);
  li.append(a);
  breadcrump.append(li);
}

function init () {
    displayLogoutButton();
    document.querySelector("#logout").addEventListener('click',logout,false);
}


document.addEventListener('DOMContentLoaded',init,false);