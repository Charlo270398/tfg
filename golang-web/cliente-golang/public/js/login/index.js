function submit(event){
    let identificacion = document.querySelector("#identificacion").value;
    let password = document.querySelector("#password").value;
    if(identificacion && password){
        login(identificacion, password);
    }
}
function login(identificacion, password){
    const url= `/login`;
    const payload= {identificacion: identificacion, password: password};
    const request = {
        method: 'POST', 
        headers: cabeceras,
        body: JSON.stringify(payload),
    };
    fetch(url,request)
    .then( response => response.json() )
        .then( r => {
            if(!r.Error){
                console.log("SESION INICIADA");
                window.location.href="/user/menu";
            }
            else{
                alert(r.Error);
            }
        })
        .catch(err => alert(err));
}

function init () {
    deleteBreadcrumb();
    addLinkBreadcrumb('Home', '/home');
    addLinkBreadcrumb('Login', '/login');
    document.querySelector("#submit").addEventListener('click',submit,false);
}

document.addEventListener('DOMContentLoaded',init,false);