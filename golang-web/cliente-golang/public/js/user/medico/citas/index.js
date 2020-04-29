function restAddEntrada(motivoConsulta, juicioDiagnostico){
    const url= `/user/doctor/historial/addEntrada`;
    const payload= {citaId: citaId, motivoConsulta: motivoConsulta, juicioDiagnostico: juicioDiagnostico};
    const request = {
        method: 'POST', 
        headers: cabeceras,
        body: JSON.stringify(payload),
    };
    fetch(url,request)
    .then( response => response.json() )
        .then( r => {
            if(!r.Error){
                //Cerrar cita
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

function addEntrada(event){
    if(citaId && document.querySelector("#motivoConsulta").value != "" && document.querySelector("#juicioDiagnostico").value != ""){
     console.log(citaId);
    }else{
        document.querySelector("#alert").textContent = "Existen campos vacíos";
        document.querySelector("#alert").classList.replace("alert-success", "alert-danger");
        document.querySelector("#alert").classList.remove('invisible');
    }
}

function init () {
    deleteBreadcrumb();
    addLinkBreadcrumb('Usuario', '/user/menu');
    addLinkBreadcrumb('Medico', '/user/doctor');
    addLinkBreadcrumb('Pasar consulta', '');
    document.querySelector("#addEntrada").addEventListener('click',addEntrada,false);
}

document.addEventListener('DOMContentLoaded',init,false);