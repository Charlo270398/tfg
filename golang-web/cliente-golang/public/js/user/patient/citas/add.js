function init () {
    loadClinicas(clinicas);
    deleteBreadcrumb();
    addLinkBreadcrumb('Usuario', '/user/menu');
    addLinkBreadcrumb('Paciente', '/user/patient');
    addLinkBreadcrumb('Citas', '/user/patient/citas');
    addLinkBreadcrumb('Solicitar', '/user/patient/citas/add');
    document.querySelector("#clinicaSelector").addEventListener('change',cambiarClinica,false);
    document.querySelector("#especialidadSelector").addEventListener('change',cambiarEspecialidad,false);
    document.querySelector("#facultativoSelector").addEventListener('change',cambiarFacultativo,false);
    document.querySelector("#diaSelector").addEventListener('change',cambiarDia,false);
    document.querySelector("#horaSelector").addEventListener('change',cambiarHora,false);
}

function loadClinicas(clinicas){
    var clinicaSelector = document.querySelector("#clinicaSelector");
    clinicas.forEach(c => {
        var option = document.createElement("option");
        option.value = c.id;
        option.textContent = c.nombre;
        clinicaSelector.append(option);
    });
}

function limpiarSelect(nodoSelect){
    var children = Array.from(nodoSelect.children);
    children.forEach(option => {if(option.value != "-1"){option.remove()}});
}

function cambiarClinica(event){
    //Recargar las especialidades
    document.querySelector("#especialidadGroup").classList.remove("invisible");
    document.querySelector("#especialidadSelector").value = "-1";
    limpiarSelect(document.querySelector("#especialidadSelector"));
    GETespecialidades(document.querySelector("#clinicaSelector").value, document.querySelector("#especialidadSelector"));

    //Recargar los facultativos
    document.querySelector("#facultativoGroup").classList.add("invisible");
    document.querySelector("#facultativoSelector").value = "-1";
    limpiarSelect(document.querySelector("#facultativoSelector"));

    //Recargar los dias
    document.querySelector("#diaGroup").classList.add("invisible");
    document.querySelector("#diaSelector").value = "-1";
    limpiarSelect(document.querySelector("#diaSelector"));

    //Recargar las horas
    document.querySelector("#horaGroup").classList.add("invisible");
    document.querySelector("#horaSelector").value = "-1";
    limpiarSelect(document.querySelector("#horaSelector"));

    
}

function cambiarEspecialidad(event){
    //Recargar los facultativos
    document.querySelector("#facultativoGroup").classList.remove("invisible");
    document.querySelector("#facultativoSelector").value = "-1";
    limpiarSelect(document.querySelector("#facultativoSelector"));
    GETfacultativos(document.querySelector("#clinicaSelector").value, 
    document.querySelector("#especialidadSelector").value, document.querySelector("#facultativoSelector"));

    //Recargar los dias
    document.querySelector("#diaGroup").classList.add("invisible");
    document.querySelector("#diaSelector").value = "-1";
    limpiarSelect(document.querySelector("#diaSelector"));
    
    //Recargar las horas
    document.querySelector("#horaGroup").classList.add("invisible");
    document.querySelector("#horaSelector").value = "-1";
    limpiarSelect(document.querySelector("#horaSelector"));
}

function cambiarFacultativo(event){
    //Recargar los dias
    document.querySelector("#diaGroup").classList.remove("invisible");
    document.querySelector("#diaSelector").value = "-1";
    limpiarSelect(document.querySelector("#diaSelector"));

    //Recargar las horas
    document.querySelector("#horaGroup").classList.add("invisible");
    document.querySelector("#horaSelector").value = "-1";
    limpiarSelect(document.querySelector("#horaSelector"));
}

function cambiarDia(event){
    //Recargar las horas
    document.querySelector("#horaGroup").classList.remove("invisible");
    document.querySelector("#horaSelector").value = "-1";
    limpiarSelect(document.querySelector("#horaSelector"));
}

function cambiarHora(event){

}

function GETespecialidades(clinica_id, selector){
    const url= `/clinica/especialidad/list?clinicaId=` + clinica_id;
    const request = {
        method: 'GET', 
        headers: cabeceras,
    };
    fetch(url,request)
    .then( response => response.json() )
        .then( result => {
            if(!result.Error){
                result.forEach(e => {
                    var option = document.createElement("option");
                    option.value = e.id;
                    option.textContent = e.nombre;
                    selector.append(option);
                });
            }
            else{

            }
        })
        .catch(err => alert(err));
}

function GETfacultativos(clinica_id, especialidad_id, selector){
    const url= `/clinica/especialidad/doctor/list?clinicaId=` + clinica_id + "&especialidadId=" + especialidad_id;
    const request = {
        method: 'GET', 
        headers: cabeceras,
    };
    fetch(url,request)
    .then( response => response.json() )
        .then( result => {
            if(!result.Error){
                result.forEach(f => {
                    var option = document.createElement("option");
                    option.value = f.Id;
                    option.textContent = f.NombreDoctor;
                    selector.append(option);
                });
            }
            else{

            }
        })
        .catch(err => alert(err));
}

document.addEventListener('DOMContentLoaded',init,false);
