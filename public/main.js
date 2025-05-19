const api = 'https://jubilant-tribble-7rwgxwp5q6p2pjjp-8081.app.github.dev';

document.addEventListener('DOMContentLoaded', () => {
    fetch(`${api}/civilisations`)
        .then(res => res.json())
        .then(civs => {
            const select = document.getElementById('civ');
            civs.forEach(civ => {
                const option = document.createElement('option');
                option.value = civ.code;
                option.textContent = civ.name;
                select.appendChild(option);
            });
            select.addEventListener('change', (e) => {
                const civCode = e.target.value;
                if (civCode) {
                    fetch(`${api}/civilisations/${civCode}/units`)
                        .then(res => res.json())
                        .then(units => {
                            // Group units by type
                            const groups = {
                                infantry: [],
                                cavalry: [],
                                champion: [],
                                hero: [],
                                ship: [],
                                siege: [],
                                support: [],
                                other: []
                            };
                            units.forEach(unit => {
                                const code = unit.code.toLowerCase();
                                if (code.includes('infantry')) groups.infantry.push(unit);
                                else if (code.includes('cavalry')) groups.cavalry.push(unit);
                                else if (code.includes('champion')) groups.champion.push(unit);
                                else if (code.includes('hero')) groups.hero.push(unit);
                                else if (code.startsWith('ship')) groups.ship.push(unit);
                                else if (code.startsWith('siege')) groups.siege.push(unit);
                                else if (code.startsWith('support')) groups.support.push(unit);
                                else groups.other.push(unit);
                            });

                            const unitsList = document.getElementById('units-list');
                            unitsList.innerHTML = '';
                            const groupOrder = [
                                { key: 'infantry', label: 'Infantry' },
                                { key: 'cavalry', label: 'Cavalry' },
                                { key: 'champion', label: 'Champions' },
                                { key: 'hero', label: 'Heroes' },
                                { key: 'ship', label: 'Ships' },
                                { key: 'siege', label: 'Siege' },
                                { key: 'support', label: 'Support' },
                                { key: 'other', label: 'Other' }
                            ];
                            groupOrder.forEach(group => {
                                if (groups[group.key].length) {
                                    const groupTitle = document.createElement('h3');
                                    groupTitle.textContent = group.label;
                                    unitsList.appendChild(groupTitle);
                                    const groupDiv = document.createElement('div');
                                    groupDiv.className = 'units-group';
                                    groupDiv.style.display = 'flex';
                                    groupDiv.style.flexWrap = 'wrap';
                                    groupDiv.style.gap = '1rem';
                                    groups[group.key].forEach(unit => {
                                        const card = document.createElement('div');
                                        card.className = 'unit-card';
                                        card.tabIndex = 0;

                                        // Add unit icon
                                        if (unit.Identity && unit.Identity.Icon) {
                                            const icon = document.createElement('img');
                                            icon.src = unit.Identity.Icon;
                                            icon.alt = unit.Identity.SpecificName || 'Unit icon';
                                            icon.style.width = '100%'; // Adjust as needed
                                            icon.style.height = 'auto'; // Adjust as needed
                                            icon.style.marginBottom = '0.5rem';
                                            card.appendChild(icon);
                                        }

                                        // Add SpecificName
                                        if (unit.Identity && unit.Identity.SpecificName) {
                                            const specificNameElement = document.createElement('div');
                                            specificNameElement.textContent = unit.Identity.SpecificName;
                                            specificNameElement.style.fontWeight = 'bold';
                                            card.appendChild(specificNameElement);
                                        }

                                        // Add GenericName
                                        if (unit.Identity && unit.Identity.GenericName) {
                                            const genericNameElement = document.createElement('div');
                                            genericNameElement.textContent = unit.Identity.GenericName;
                                            genericNameElement.style.fontSize = '0.8em';
                                            genericNameElement.style.color = '#666';
                                            card.appendChild(genericNameElement);
                                        }
                                        // TODO: Add click handler to show unit details
                                        groupDiv.appendChild(card);
                                    });
                                    unitsList.appendChild(groupDiv);
                                }
                            });
                        })
                        .catch(err => {
                            console.error('Failed to load units:', err);
                        });
                } else {
                    // TODO: Clear units-list section
                }
            });
        })
        .catch(err => {
            console.error('Failed to load civilisations:', err);
        });
});