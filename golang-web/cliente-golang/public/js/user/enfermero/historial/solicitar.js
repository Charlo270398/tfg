function busquedaNombre(event){
    document.querySelector("#busquedaLabel").textContent = "Nombre";
}

function busquedaDni(event){
    document.querySelector("#busquedaLabel").textContent = "Documento de identificación";
}


function init () {
    deleteBreadcrumb();
    addLinkBreadcrumb('Usuario', '/user/menu');
    addLinkBreadcrumb('Enfermero', '/user/nurse');
    addLinkBreadcrumb('Solicitar historial', '/user/historial/solicitar');
    document.querySelector("#radioNombre").addEventListener('change',busquedaNombre,false);
    document.querySelector("#radioDni").addEventListener('change',busquedaDni,false);
}

document.addEventListener('DOMContentLoaded',init,false);