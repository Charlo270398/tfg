function init () {
    loadClinicas(clinicas);
    deleteBreadcrumb();
    addLinkBreadcrumb('Usuario', '/user/menu');
    addLinkBreadcrumb('Paciente', '/user/patient');
    addLinkBreadcrumb('Citas', '/user/patient/citas');
    addLinkBreadcrumb('Solicitar', '/user/patient/citas/add');
    document.querySelector("#clinicaSelector").addEventListener('change',cambiarClinica,false);
    document.querySelector("#especialidadSelector").addEventListener('change',cambiarEspecialidad,false);
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

document.addEventListener('DOMContentLoaded',init,false);
