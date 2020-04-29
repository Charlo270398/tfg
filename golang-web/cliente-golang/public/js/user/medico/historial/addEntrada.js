function submit(event){
    console.log(document.querySelector("#inputFile").value);   
}

function cargarTablaHistorial(historial){
    document.querySelector("#alert").classList.add('invisible');
    document.querySelector("#historialTabla").classList.remove('invisible');
}

function restBuscarDNI(DNI){
    const url= `/user/doctor/historial/solicitar`;
    const payload= {identificacion: DNI};
    const request = {
        method: 'POST', 
        headers: cabeceras,
        body: JSON.stringify(payload),
    };
    fetch(url,request)
    .then( response => response.json() )
        .then( r => {
            if(!r.Error){
                //PROCESAR HISTORIAL
                console.log(r);
                cargarTablaHistorial(r);
            }
            else{
                document.querySelector("#alert").textContent = r.Error;
                document.querySelector("#alert").classList.replace("alert-success", "alert-danger");
                document.querySelector("#alert").classList.remove('invisible');
                document.querySelector("#historialTabla").classList.add('invisible');
            }
        })
        .catch(err => alert(err));
}

function init () {
    deleteBreadcrumb();
    addLinkBreadcrumb('Usuario', '/user/menu');
    addLinkBreadcrumb('Medico', '/user/doctor');
    addLinkBreadcrumb('Consultar historiales', '/user/historial/consultar');
    document.querySelector("#searchButton").addEventListener('click',submit,false);
}

document.addEventListener('DOMContentLoaded',init,false);