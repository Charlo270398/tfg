(function() {
    const cabeceras= {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
    }
    const template = document.createElement('template');
    template.innerHTML = `<label><strong>Tags de la analítica</strong></label>   
                          <div id="tagsGroup" class="checkbox"></div>
                          `;

    function addCheckBox(groupDiv, tag){
        //Añadir las cajas con roles dinamicamente
        let div = document.createElement('div');
        div.classList.add('checkbox');
        div.setAttribute('tag_id',tag.id);
        let lab = document.createElement('label');
        let input = document.createElement('input');
        input.type = 'checkbox';
        lab.append(input);
        lab.append(tag.nombre);
        div.append(lab);
        groupDiv.append(div);
        input.addEventListener('change',changeCheckBox,false);
        return input;
    }
    
    function cargarTags(groupDiv, select){
        const url= `/tag/list`;
        const request = {
            method: 'GET', 
            headers: cabeceras,
        };
        fetch(url,request)
        .then( response => response.json() )
            .then( r => {
                if(!r.Error){
                    r.forEach(tag => {
                        checkbox = addCheckBox(groupDiv.querySelector("#tagsGroup"), tag);
                        if(select){
                            
                        }
                    });
                }
                else{
                    console.log("ERROR CARGANDO TAGS");
                }
            })
            .catch(err => alert(err));
    }

    function changeCheckBox(event){
        if(event.target.checked){
            tagsArray.push(parseInt(event.target.parentNode.parentNode.getAttribute('tag_id')));
        }else{
            var pos = tagsArray.indexOf(parseInt(event.target.parentNode.parentNode.getAttribute('tag_id')));
            tagsArray.splice(pos, 1);
        }
    }
  
    class Checkbox_tags extends HTMLElement {
        
        constructor() {
            super();
            let tclone = template.content.cloneNode(true);
            let shadowRoot = this.attachShadow({
                mode: 'open' 
            });
            shadowRoot.appendChild(tclone);
        }
  
        connectedCallback() {
            cargarTags(this.shadowRoot, true);
        }
    }
  
    customElements.define("checkbox-tags", Checkbox_tags); //Definimos el nombre del componente
  
  })();