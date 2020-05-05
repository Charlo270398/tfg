function submit(event){
    let nombre = document.querySelector("#name").value;
    let apellido1 = document.querySelector("#surname1").value;
    let apellido2 = document.querySelector("#surname2").value;
    let identificacion = document.querySelector("#idnumber").value;
    let email = document.querySelector("#email").value;
    if(nombre && apellido1 && apellido2 && identificacion && email){
        register(nombre,apellido1, apellido2, email, identificacion);
    }else{
        document.querySelector("#alert").textContent = "Existen campos vacíos";
        document.querySelector("#alert").classList.replace("alert-success", "alert-danger");
        document.querySelector("#alert").classList.remove('invisible');
    }
}
function register(nombre, apellido1, apellido2, email, identificacion){
    var result = false;
    let apellidos = apellido1;
    if(apellido2){
        apellidos += " " + apellido2;
    }
    const url= `/user/menu/edit`;
    const payload= {nombre: nombre, identificacion:identificacion, apellidos: apellidos, email: email};
    const request = {
        method: 'POST', 
        headers: cabeceras,
        body: JSON.stringify(payload),
    };
    fetch(url,request)
    .then( response => response.json() )
        .then( r => {
            if(!r.Error){
                console.log("USUARIO MODIFICADO CORRECTAMENTE");
                document.querySelector("#alert").classList.replace("alert-danger", "alert-success");
                document.querySelector("#alert").textContent = "Datos del usuario modificados correctamente";
                document.querySelector("#alert").classList.remove('invisible');
                //window.location.href="/user/menu";
            }
            else{
                document.querySelector("#alert").textContent = r.Error;
                document.querySelector("#alert").classList.replace("alert-success", "alert-danger");
                document.querySelector("#alert").classList.remove('invisible');
            }
        })
        .catch(err => alert(err));
    return result;
}
function cargarDatos(){
    document.querySelector("#name").value = nombre;
    document.querySelector("#surname1").value = apellido1;
    document.querySelector("#surname2").value = apellido2;
    document.querySelector("#email").value = email;
}
function init () {
    deleteBreadcrumb();
    addLinkBreadcrumb('Usuario', '/user/menu');
    addLinkBreadcrumb('Editar', '/user/menu/edit');
    cargarDatos(nombre, apellido1, apellido2);
    document.querySelector("#submit").addEventListener('click',submit,false);
}

document.addEventListener('DOMContentLoaded',init,false);