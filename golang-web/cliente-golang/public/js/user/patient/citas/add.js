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
    document.querySelector("#horaSelector").addEventListener('change',cambiarHora,false);
    document.querySelector("#submit").addEventListener('click',submit,false);
}

function submit(event){
    let facultativoId = document.querySelector("#facultativoSelector").value;
    let fecha = document.querySelector("#horaSelector").value;
    if(facultativoId && fecha){
        registrarCita(fecha, facultativoId);
    }else{
        document.querySelector("#alert").textContent = "Existen campos vacíos";
        document.querySelector("#alert").classList.replace("alert-success", "alert-danger");
        document.querySelector("#alert").classList.remove('invisible');
        return;
    }
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

    document.querySelector("#submit").classList.add("invisible");
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

    document.querySelector("#submit").classList.add("invisible");
}

function cambiarFacultativo(event){
    //Recargar los dias
    document.querySelector("#diaGroup").classList.remove("invisible");
    document.querySelector("#diaSelector").value = "-1";
    limpiarSelect(document.querySelector("#diaSelector"));
    GETdias(document.querySelector("#facultativoSelector").value,
    document.querySelector("#diaSelector"));

    //Recargar las horas
    document.querySelector("#horaGroup").classList.add("invisible");
    document.querySelector("#horaSelector").value = "-1";
    limpiarSelect(document.querySelector("#horaSelector"));

    document.querySelector("#submit").classList.add("invisible");
}

function cambiarDia(event){
    //Recargar las horas
    document.querySelector("#horaGroup").classList.remove("invisible");
    document.querySelector("#horaSelector").value = "-1";
    limpiarSelect(document.querySelector("#horaSelector"));
    GEThoras(document.querySelector("#facultativoSelector").value,
    document.querySelector("#diaSelector").value,
    document.querySelector("#horaSelector"));

    document.querySelector("#submit").classList.add("invisible");
}

function cambiarHora(event){
    document.querySelector("#submit").classList.remove("invisible");
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
            if(result){
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
                    option.value = f.id;
                    option.textContent = f.nombreDoctor;
                    selector.append(option);
                });
            }
            else{

            }
        })
        .catch(err => alert(err));
}

function GETdias(facultativo_id, selector){
    const url= `/user/doctor/disponible/dia?doctorId=` + facultativo_id;
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
                    //Tratamos para usar en el servidor
                    var diaTratado = f.Dia
                    var  mesTratado = f.Mes
                    if(parseInt(f.Dia)<10){
                        diaTratado = "0"+f.Dia
                    }
                    if(parseInt(f.Mes)<10){
                        mesTratado = "0"+f.Mes
                    }
                    option.value = f.Anyo + "-" + mesTratado + "-" + diaTratado;
                    option.textContent = f.Dia + "-" + f.Mes + "-" + f.Anyo;
                    selector.append(option);
                });
            }
            else{

            }
        })
        .catch(err => alert(err));
}

function GEThoras(facultativo_id, dia, selector){
    const url= `/user/doctor/disponible/hora?doctorId=` + facultativo_id + "&dia=" + dia;
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
                    //Tratamos para usar en el servidor
                    var diaTratado = f.Dia
                    var  mesTratado = f.Mes
                    if(parseInt(f.Dia)<10){
                        diaTratado = "0"+f.Dia
                    }
                    if(parseInt(f.Mes)<10){
                        mesTratado = "0"+f.Mes
                    }
                    option.value = f.Fecha.substring(0, f.Fecha.length - 1);
                    option.textContent = f.Hora + ":00";
                    selector.append(option);
                });
            }
            else{

            }
        })
        .catch(err => alert(err));
}

function registrarCita(fecha, facultativoId){
    const url= `/user/patient/citas/add`;
    const payload= {medicoId: facultativoId, fechaString: fecha};
    const request = {
        method: 'POST', 
        headers: cabeceras,
        body: JSON.stringify(payload),
    };
    fetch(url,request)
    .then( response => response.json() )
        .then( r => {
            if(!r.Error){
                console.log("CITA AÑADIDA CORRECTAMENTE");
                document.querySelector("#submit").classList.add("invisible");
                //Recargar las clinica
                document.querySelector("#clinicaSelector").value = "-1";

                //Recargar las especialidades
                document.querySelector("#especialidadGroup").classList.add("invisible");
                document.querySelector("#especialidadSelector").value = "-1";
                limpiarSelect(document.querySelector("#especialidadSelector"));

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
                
                document.querySelector("#alert").classList.replace("alert-danger", "alert-success");
                document.querySelector("#alert").textContent = "Cita reservada correctamente";
                document.querySelector("#alert").classList.remove('invisible');
            }
            else{
                document.querySelector("#alert").textContent = r.Error;
                document.querySelector("#alert").classList.replace("alert-success", "alert-danger");
                document.querySelector("#alert").classList.remove('invisible');
            }
        })
        .catch(err => alert(err));
}

document.addEventListener('DOMContentLoaded',init,false);
