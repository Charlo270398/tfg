function submit(event){
    let nombre = document.querySelector("#name").value;
    if(nombre){
        insertarEspecialidad(nombre);
    }else{
        document.querySelector("#alert").textContent = "Existen campos vacíos";
        document.querySelector("#alert").classList.replace("alert-success", "alert-danger");
        document.querySelector("#alert").classList.remove('invisible');
    }
}
function insertarEspecialidad(nombre){
    const url= `/especialidad/add`;
    const payload= {nombre: nombre};
    const request = {
        method: 'POST', 
        headers: cabeceras,
        body: JSON.stringify(payload),
    };
    fetch(url,request)
    .then( response => response.json() )
        .then( r => {
            if(!r.Error){
                console.log("ESPECIALIDAD AÑADIDA CORRECTAMENTE");
                document.querySelector("#alert").classList.replace("alert-danger", "alert-success");
                document.querySelector("#alert").textContent = "Especialidad añadida correctamente";
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
    addLinkBreadcrumb('Especialidades', '/especialidad/list');
    addLinkBreadcrumb('Añadir', '/especialidad/add');
    document.querySelector("#submit").addEventListener('click',submit,false);
}

document.addEventListener('DOMContentLoaded',init,false);