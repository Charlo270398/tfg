function submit(event){
    let nombre = document.querySelector("#name").value;
    let apellido1 = document.querySelector("#surname1").value;
    let apellido2 = document.querySelector("#surname2").value;
    let identificacion = document.querySelector("#idnumber").value;
    let email = document.querySelector("#email").value;
    let password = document.querySelector("#pass").value;
    let condiciones = document.querySelector("#conditions").checked;
    if(nombre && apellido1 && apellido2 && identificacion && email && password && condiciones){
        register(nombre,apellido1, apellido2, email, identificacion, password);
    }
}
function register(nombre, apellido1, apellido2, email, identificacion, password){
    var result = false;
    let apellidos = apellido1;
    if(apellido2){
        apellidos += " " + apellido2;
    }
    const url= `/register`;
    const payload= {nombre: nombre, identificacion:identificacion, apellidos: apellidos, email: email, password: password};
    const request = {
        method: 'POST', 
        headers: cabeceras,
        body: JSON.stringify(payload),
    };
    fetch(url,request)
    .then( response => response.json() )
        .then( r => {
            if(!r.Error){
                //Registrado, continuamos a menú
                console.log("USUARIO REGISTRADO CORRECTAMENTE");
                window.location.href="/user/menu";
            }
            else{
                alert(r.Error);
            }
        })
        .catch(err => alert(err));
    return result;
}
function init () {
    deleteBreadcrumb();
    addLinkBreadcrumb('Home', '/home');
    addLinkBreadcrumb('Register', '/register');
    document.querySelector("#submit").addEventListener('click',submit,false);
}

document.addEventListener('DOMContentLoaded',init,false);