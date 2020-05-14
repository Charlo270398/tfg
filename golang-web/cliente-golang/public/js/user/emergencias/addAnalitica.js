function submit(event){
    
}

function restAddAnalitica(DNI){
    const url= `/user/doctor/historial/addAnalitica`;
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

function changeNumber(event){
    if(event.target.value >= event.target.max){
        event.target.value = event.target.max;
    }
    console.log(parseFloat(event.target.value));
}

function init () {
    deleteBreadcrumb();
    addLinkBreadcrumb('Usuario', '/user/menu');
    addLinkBreadcrumb('Emergencias', '/user/emergency');
    //Si se pasa por parametro el DNI se busca auto
    var url = new URL(window.location.href);
    var paramIdentificacion = url.searchParams.get("identificacion");
    if(paramIdentificacion){
        addLinkBreadcrumb('Historial', '/user/emergency?identificacion='+paramIdentificacion);
    }
    addLinkBreadcrumb('Añadir analítica', '');
    document.querySelector("#submit").addEventListener('click',submit,false);

    document.querySelector("#numberLeucocitos").addEventListener('change',changeNumber,false);
}

document.addEventListener('DOMContentLoaded',init,false);