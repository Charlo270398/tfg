function restAddEntrada(motivoConsulta, juicioDiagnostico){
    const url= `/user/doctor/citas/addEntrada`;
    const payload= {citaId: parseInt(cita.id), pacienteId: parseInt(cita.pacienteId), motivoConsulta: motivoConsulta, juicioDiagnostico: juicioDiagnostico};
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
                console.log(r);
                if(r.Result == "OK"){
                    document.querySelector("#alert").textContent = "Entrada insertada correctamente";
                    document.querySelector("#alert").classList.remove('invisible');
                    document.querySelector("#formConsulta").classList.add('invisible');
                }
            }
            else{
                document.querySelector("#alert").textContent = r.Error;
                document.querySelector("#alert").classList.replace("alert-success", "alert-danger");
                document.querySelector("#alert").classList.remove('invisible');
            }
        })
        .catch(err => alert(err));
}

function addEntrada(event){
    if(cita.id && document.querySelector("#motivoConsulta").value != "" && document.querySelector("#juicioDiagnostico").value != ""){
     restAddEntrada(document.querySelector("#motivoConsulta").value, document.querySelector("#juicioDiagnostico").value);
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