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
                const unitsList = document.getElementById('units-list');
                const spinner = document.getElementById('loading-spinner');

                if (civCode) {
                    if (spinner) {
                        spinner.style.display = 'block'; // Show spinner
                    } else {
                        console.error("Spinner element #loading-spinner not found when trying to show it.");
                    }
                    unitsList.innerHTML = ''; // Clear previous units
                    // Ensure the spinner is the only child if unitsList was cleared and spinner exists
                    if (spinner && !unitsList.contains(spinner)) {
                        unitsList.appendChild(spinner); 
                    }

                    fetch(`${api}/civilisations/${civCode}/units`)
                        .then(res => res.json())
                        .then(units => {
                            // START: Filter units to display one per SpecificName and mark if others exist
                            const specificNameCounts = units.reduce((acc, unit) => {
                                const name = unit.Identity?.SpecificName;
                                if (name) {
                                    acc[name] = (acc[name] || 0) + 1;
                                }
                                return acc;
                            }, {});

                            const uniqueDisplayUnits = [];
                            const seenSpecificNames = new Set();

                            units.forEach(originalUnit => {
                                const unit = { ...originalUnit }; // Create a copy to safely add properties
                                const specificName = unit.Identity?.SpecificName;

                                if (specificName) {
                                    // This unit has a SpecificName
                                    if (!seenSpecificNames.has(specificName)) {
                                        // This is the first time we're seeing this SpecificName
                                        if (specificNameCounts[specificName] > 1) {
                                            unit.hasMultipleVariants = true; // Mark if other variants exist
                                        }
                                        uniqueDisplayUnits.push(unit);
                                        seenSpecificNames.add(specificName);
                                    }
                                    // If seenSpecificNames.has(specificName) is true, we skip this unit.
                                } else {
                                    // Unit does not have a SpecificName, display it as is.
                                    uniqueDisplayUnits.push(unit);
                                }
                            });
                            // END: Filter units

                            // Group units by type using uniqueDisplayUnits
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
                            uniqueDisplayUnits.forEach(unit => { // MODIFIED: Use uniqueDisplayUnits here
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

                            // Hide the spinner now that data is processed and before adding new elements
                            if (spinner) {
                                spinner.style.display = 'none'; 
                            }

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

                                            if (unit.hasMultipleVariants) {
                                                const multiIcon = document.createElement('span');
                                                multiIcon.textContent = ' ➕'; // Using a plus icon, with a leading space
                                                multiIcon.title = 'Multiple variants or upgrades exist for this unit type';
                                                multiIcon.style.fontSize = '0.9em'; // Make icon slightly smaller
                                                multiIcon.style.fontWeight = 'normal'; // Keep icon normal weight
                                                specificNameElement.appendChild(multiIcon);
                                            }
                                            card.appendChild(specificNameElement);
                                        }

                                        // Add GenericName
                                        if (unit.Identity && (unit.Identity.GenericName || unit.Identity.VisibleClasses)) {
                                            const genericNameElement = document.createElement('div');
                                            genericNameElement.textContent = unit.Identity.VisibleClasses ?? unit.Identity.GenericName;
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
                            if (spinner) {
                                spinner.style.display = 'none'; // Hide spinner on error
                            }
                        });
                } else {
                    unitsList.innerHTML = ''; // Clear units-list section
                    if (spinner) {
                        spinner.style.display = 'none'; // Hide spinner if no civ is selected
                    }
                }
            });
        })
        .catch(err => {
            console.error('Failed to load civilisations:', err);
        });
});