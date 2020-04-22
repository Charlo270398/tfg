function submit(event){
    let email = document.querySelector("#email").value;
    let password = document.querySelector("#password").value;
    if(email && password){
        login(email, password);
    }
}
function login(email, password){
    const url= `/login`;
    const payload= {email: email, password: password};
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