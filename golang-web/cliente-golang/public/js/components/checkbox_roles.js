(function() {
    const cabeceras= {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
    }
    const template = document.createElement('template');
    template.innerHTML = `<label><strong>Roles del usuario</strong></label>   
                          <div id="rolesGroup" class="checkbox"></div>
                          `;

    function addCheckBox(groupDiv, rol){
        //Añadir las cajas con roles dinamicamente
        let div = document.createElement('div');
        div.classList.add('checkbox');
        div.setAttribute('rol_id',rol.Id);
        let lab = document.createElement('label');
        let input = document.createElement('input');
        input.type = 'checkbox';
        lab.append(input);
        lab.append(rol.Descripcion);
        div.append(lab);
        groupDiv.append(div);
        return input;
    }
    
    function cargarRoles(groupDiv, select){
        const url= `/rol/list`;
        const request = {
            method: 'GET', 
            headers: cabeceras,
        };
        fetch(url,request)
        .then( response => response.json() )
            .then( r => {
                if(!r.Error){
                    r.roles.forEach(rol => {
                        checkbox = addCheckBox(groupDiv.querySelector("#rolesGroup"), rol);
                        if(select){
                            if(rol.Id == "2"){
                                checkbox.addEventListener('change',mostrarEnfermeroClinicaGroup,false); 
                            }
                            if(rol.Id == "3"){
                                checkbox.addEventListener('change',mostrarMedicoClinicaGroup,false); 
                            }
                            if(rol.Id == "4"){
                                checkbox.addEventListener('change',mostrarAdminClinicaGroup,false); 
                            }
                        }
                    });
                    //window.location.href="/user/menu";
                }
                else{
                    console.log("ERROR CARGANDO ROLES");
                }
            })
            .catch(err => alert(err));
    }
  
    class Checkbox_roles extends HTMLElement {
        constructor() {
            super();
            let tclone = template.content.cloneNode(true);
            let shadowRoot = this.attachShadow({
                mode: 'open' 
            });
            shadowRoot.appendChild(tclone);
        }
  
        connectedCallback() {
            if(this.hasAttribute('userId')){//Si se introduce userId como atributo

            }else{//Si está vacío
                cargarRoles(this.shadowRoot, true);
            }
        }
    }
  
    customElements.define("checkbox-roles", Checkbox_roles); //Definimos el nombre del componente
  
  })();