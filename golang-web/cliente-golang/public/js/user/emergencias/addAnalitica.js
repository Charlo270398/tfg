const tagsArray = [];
var HISTORIAL_ID;

function submit(event){
    if(tagsArray.length <= 0){
        document.querySelector("#alert").textContent = "Debes añadir al menos un tag";
        document.querySelector("#alert").classList.replace("alert-success", "alert-danger");
        document.querySelector("#alert").classList.remove('invisible');
    }else{
        if(HISTORIAL_ID != undefined){
            document.querySelector("#alert").classList.add('invisible');
            const numberLeucocitos = document.querySelector("#numberLeucocitos").value;
            const numberHematies = document.querySelector("#numberHematies").value;
            const numberPlaquetas = document.querySelector("#numberPlaquetas").value;
            const numberGlucosa = document.querySelector("#numberGlucosa").value;
            const numberHierro = document.querySelector("#numberHierro").value;
            restAddAnalitica(numberLeucocitos, numberHematies, numberPlaquetas, numberGlucosa, numberHierro, tagsArray);
        }
    }
}

function restAddAnalitica(numberLeucocitos, numberHematies, numberPlaquetas, numberGlucosa, numberHierro, tags){
    const url= `/user/emergency/historial/addAnalitica`;
    const payload= {historialId: HISTORIAL_ID, leucocitos: numberLeucocitos, hematies: numberHematies, plaquetas: numberPlaquetas, glucosa: numberGlucosa, hierro: numberHierro, tags: tags};
    const request = {
        method: 'POST', 
        headers: cabeceras,
        body: JSON.stringify(payload),
    };
    fetch(url,request)
    .then( response => response.json() )
        .then( r => {
            if(!r.Error){
                if(r.Result == "OK"){
                    document.querySelector("#alert").textContent = "Analítica insertada correctamente";
                    document.querySelector("#alert").classList.remove('invisible');
                    document.querySelector("#formConsulta").classList.add('invisible');
                }
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
    if(event.target.value < event.target.min){
        event.target.value = 0;
    }
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
    var paramHistorialId = url.searchParams.get("historialId");
    if(paramHistorialId){
        HISTORIAL_ID = parseInt(paramHistorialId);
    }
    //Eventos
    document.querySelector("#submit").addEventListener('click',submit,false);
    const datos = document.querySelector("#datos");
    inputsDatos = datos.querySelectorAll('input');
    inputsDatos.forEach((input) => {
        input.addEventListener('change',changeNumber,false);
    });
}

document.addEventListener('DOMContentLoaded',init,false);