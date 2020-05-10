var HISTORIAL_ID = -1;

function busquedaDNI(event){
    if(document.querySelector("#inputDNI").value.length != 9){
        //Activar alerta
        document.querySelector("#alert").textContent = "El documento de identificación debe tener un formato válido (por ejemplo, 00000000X)";
        document.querySelector("#alert").classList.replace("alert-success", "alert-danger");
        document.querySelector("#alert").classList.remove('invisible');
        document.querySelector("#historialDiv").classList.add('invisible');
        return;
    }else{
        restBuscarDNI(document.querySelector("#inputDNI").value);
    }
}

function cargarTablaHistorial(entradas){
    console.log(entradas);
    document.querySelector("#alert").classList.add('invisible');
    document.querySelector("#historialTabla").classList.remove('invisible');
}

function restBuscarDNI(DNI){
    const url= `/user/emergency/historial`;
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
                if(r.id != 0){
                    //PROCESAR HISTORIAL
                    HISTORIAL_ID = r.id;
                    console.log(r);
                    document.querySelector("#alert").classList.add('invisible');
                    document.querySelector("#historialDiv").classList.remove('invisible');
                    document.querySelector("#spanNombre").textContent = r.nombrePaciente + " " + r.apellidosPaciente;
                    document.querySelector("#spanSexo").textContent = r.sexo;
                    document.querySelector("#spanAlergias").textContent = r.alergias
                    if(r.entradas == null){
                        document.querySelector("#alertTablaHistorial").classList.remove('invisible');
                    }else{
                        cargarTablaHistorial(r.entradas);
                    }
                }else{
                    document.querySelector("#historialDiv").classList.add('invisible');
                    document.querySelector("#alert").textContent = "No existe ningún usuario con esa identificación";
                    document.querySelector("#alert").classList.replace("alert-success", "alert-danger");
                    document.querySelector("#alert").classList.remove('invisible');
                }
            }
            else{
                document.querySelector("#alert").textContent = r.Error;
                document.querySelector("#alert").classList.replace("alert-success", "alert-danger");
                document.querySelector("#alert").classList.remove('invisible');
                document.querySelector("#historialDiv").classList.add('invisible');
            }
        })
        .catch(err => alert(err));
}

function addEntrada(){
    window.location.href = "/user/emergency/historial/addEntrada?historialId=" + HISTORIAL_ID;
}

function addAnalitica(){
    window.location.href = "/user/emergency/historial/addAnalitica?historialId=" + HISTORIAL_ID;
}


function init () {
    deleteBreadcrumb();
    addLinkBreadcrumb('Usuario', '/user/menu');
    addLinkBreadcrumb('Emergencias', '/user/emergency');
    document.querySelector("#searchButton").addEventListener('click',busquedaDNI,false);
    document.querySelector("#addEntradaButton").addEventListener('click',addEntrada,false);
    document.querySelector("#addAnaliticaButton").addEventListener('click',addAnalitica,false);
}

document.addEventListener('DOMContentLoaded',init,false);