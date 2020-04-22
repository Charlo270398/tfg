function submit(event){
    let nombre = document.querySelector("#name").value;
    let direccion = document.querySelector("#direccion").value;
    let telefono = document.querySelector("#telefono").value;
    if(nombre && direccion && telefono){
        insertarClinica(nombre, direccion, telefono);
    }else{
        document.querySelector("#alert").textContent = "Existen campos vacíos";
        document.querySelector("#alert").classList.replace("alert-success", "alert-danger");
        document.querySelector("#alert").classList.remove('invisible');
    }
}
function insertarClinica(nombre, direccion, telefono){
    const url= `/clinica/add`;
    const payload= {nombre: nombre, direccion: direccion, telefono: telefono};
    const request = {
        method: 'POST', 
        headers: cabeceras,
        body: JSON.stringify(payload),
    };
    fetch(url,request)
    .then( response => response.json() )
        .then( r => {
            if(!r.Error){
                console.log("CLINICA AÑADIDA CORRECTAMENTE");
                document.querySelector("#alert").classList.replace("alert-danger", "alert-success");
                document.querySelector("#alert").textContent = "Clínica añadida correctamente";
                document.querySelector("#alert").classList.remove('invisible');
            }
            else{
                document.querySelector("#alert").textContent = r.Error;
                document.querySelector("#alert").classList.replace("alert-success", "alert-danger");
                document.querySelector("#alert").classList.remove('invisible');
            }
        })
        .catch(err => document.querySelector("#alert").textContent = err,
        document.querySelector("#alert").classList.replace("alert-success", "alert-danger"),
        document.querySelector("#alert").classList.remove('invisible'));
    return result;
}

function init () {
    deleteBreadcrumb();
    addLinkBreadcrumb('Usuario', '/user/menu');
    addLinkBreadcrumb('Administrador global', '/user/adminG');
    addLinkBreadcrumb('Clinicas', '/clinica/list');
    addLinkBreadcrumb('Añadir especialidad', '/clinica/especialidad/add');
    document.querySelector("#submit").addEventListener('click',submit,false);
}

document.addEventListener('DOMContentLoaded',init,false);