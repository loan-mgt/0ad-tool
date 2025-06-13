const api = "";


// document.addEventListener("DOMContentLoaded", () => {
//   fetch(`${api}/civilisations`)
//     .then((res) => res.json())
//     .then((_) => {

//       if (savedCiv) {
//         // Instead of dispatching change event, call the handler directly
//         const unitsList = document.getElementById("units-list");
//         const spinner = document.getElementById("loading-spinner");
//         if (savedCiv) {
//           if (spinner) {
//             spinner.style.display = "block";
//           }
//           unitsList.innerHTML = "";
//           if (spinner && !unitsList.contains(spinner)) {
//             unitsList.appendChild(spinner);
//           }
//           fetch(`${api}/civilisations/${savedCiv}/units`)
//             .then((res) => res.json())
//             .then((units) => {
//               // START: Filter units to display one per SpecificName and mark if others exist
//               const specificNameCounts = units.reduce((acc, unit) => {
//                 const name = unit.Identity?.SpecificName;
//                 if (name) {
//                   acc[name] = (acc[name] || 0) + 1;
//                 }
//                 return acc;
//               }, {});

//               const uniqueDisplayUnits = [];
//               const seenSpecificNames = new Set();

//               units.forEach((originalUnit) => {
//                 const unit = { ...originalUnit };
//                 const specificName = unit.Identity?.SpecificName;
//                 if (specificName) {
//                   if (!seenSpecificNames.has(specificName)) {
//                     if (specificNameCounts[specificName] > 1) {
//                       unit.hasMultipleVariants = true;
//                     }
//                     uniqueDisplayUnits.push(unit);
//                     seenSpecificNames.add(specificName);
//                   }
//                 } else {
//                   uniqueDisplayUnits.push(unit);
//                 }
//               });
//               window.allUnitsForCiv = units;
//               const groups = {
//                 infantry: [],
//                 cavalry: [],
//                 champion: [],
//                 hero: [],
//                 ship: [],
//                 siege: [],
//                 support: [],
//                 other: [],
//               };
//               uniqueDisplayUnits.forEach((unit) => {
//                 const code = unit.code.toLowerCase();
//                 if (code.includes("infantry")) groups.infantry.push(unit);
//                 else if (code.includes("cavalry")) groups.cavalry.push(unit);
//                 else if (code.includes("champion")) groups.champion.push(unit);
//                 else if (code.includes("hero")) groups.hero.push(unit);
//                 else if (code.startsWith("ship")) groups.ship.push(unit);
//                 else if (code.startsWith("siege")) groups.siege.push(unit);
//                 else if (code.startsWith("support")) groups.support.push(unit);
//                 else groups.other.push(unit);
//               });
//               if (spinner) spinner.style.display = "none";
//               const groupOrder = [
//                 { key: "infantry", label: "Infantry" },
//                 { key: "cavalry", label: "Cavalry" },
//                 { key: "champion", label: "Champions" },
//                 { key: "hero", label: "Heroes" },
//                 { key: "ship", label: "Ships" },
//                 { key: "siege", label: "Siege" },
//                 { key: "support", label: "Support" },
//                 { key: "other", label: "Other" },
//               ];
//               groupOrder.forEach((group) => {
//                 if (groups[group.key].length) {
//                   const groupTitle = document.createElement("h3");
//                   groupTitle.textContent = group.label;
//                   unitsList.appendChild(groupTitle);
//                   const groupDiv = document.createElement("div");
//                   groupDiv.className = "units-group";
//                   groupDiv.style.display = "flex";
//                   groupDiv.style.flexWrap = "wrap";
//                   groupDiv.style.gap = "1rem";
//                   groups[group.key].forEach((unit) => {
//                     const card = document.createElement("div");
//                     card.className = "unit-card";
//                     card.tabIndex = 0;
//                     if (unit.Identity && unit.Identity.Icon) {
//                       const icon = document.createElement("img");
//                       icon.src = `static/${unit.Identity.Icon}`;
//                       icon.alt = unit.Identity.SpecificName || "Unit icon";
//                       icon.style.width = "100%";
//                       icon.style.height = "auto";
//                       icon.style.marginBottom = "0.5rem";
//                       card.appendChild(icon);
//                     }
//                     if (unit.Identity && unit.Identity.SpecificName) {
//                       const specificNameElement = document.createElement("div");
//                       specificNameElement.textContent =
//                         unit.Identity.SpecificName;
//                       specificNameElement.style.fontWeight = "bold";
//                       if (unit.hasMultipleVariants) {
//                         const multiIcon = document.createElement("span");
//                         multiIcon.textContent = " ➕";
//                         multiIcon.title =
//                           "Multiple variants or upgrades exist for this unit type";
//                         multiIcon.style.fontSize = "0.9em";
//                         multiIcon.style.fontWeight = "normal";
//                         specificNameElement.appendChild(multiIcon);
//                       }
//                       card.appendChild(specificNameElement);
//                     }
//                     if (
//                       unit.Identity &&
//                       (unit.Identity.GenericName ||
//                         unit.Identity.VisibleClasses)
//                     ) {
//                       const genericNameElement = document.createElement("div");
//                       genericNameElement.textContent =
//                         unit.Identity.VisibleClasses ??
//                         unit.Identity.GenericName;
//                       genericNameElement.style.fontSize = "0.8em";
//                       genericNameElement.style.color = "#666";
//                       card.appendChild(genericNameElement);
//                     }
//                     card.addEventListener("click", () => {
//                       displayUnitDetails(unit, window.allUnitsForCiv);
//                     });
//                     card.addEventListener("keydown", (e) => {
//                       if (e.key === "Enter" || e.key === " ") {
//                         displayUnitDetails(unit, window.allUnitsForCiv);
//                       }
//                     });
//                     groupDiv.appendChild(card);
//                   });
//                   unitsList.appendChild(groupDiv);
//                 }
//               });
//             })
//             .catch((err) => {
//               console.error("Failed to load units:", err);
//               if (spinner) {
//                 spinner.style.display = "none";
//               }
//             });
//         }
//       }
//     })
//     .catch((err) => {
//       console.error("Failed to load civilisations:", err);
//     });
// });

function displayUnitDetails(clickedUnit, allCivUnits) {
  const detailsSection = document.getElementById("unit-details");
  if (!detailsSection) {
    console.error("Unit details section not found.");
    return;
  }

  let unitToShow = clickedUnit;
  // If the clicked unit has multiple variants and no specific rank (or we want to default to 'Basic'),
  // find the 'Basic' rank version of this unit.
  if (clickedUnit.hasMultipleVariants) {
    const basicVariant = allCivUnits.find(
      (u) =>
        u.Identity?.SpecificName === clickedUnit.Identity?.SpecificName &&
        u.Identity?.Rank === "Basic",
    );
    if (basicVariant) {
      unitToShow = basicVariant;
    } else {
      // If no basic variant, just show the clicked one (it might be the only one or a non-basic default)
      console.warn(
        `No 'Basic' rank variant found for ${clickedUnit.Identity?.SpecificName}. Displaying clicked unit.`,
      );
    }
  }

  detailsSection.innerHTML = ""; // Clear previous details

  const title = document.createElement("h2");
  title.textContent = unitToShow.Identity?.SpecificName || "Unit Details";
  detailsSection.appendChild(title);

  if (unitToShow.Identity?.GenericName) {
    const genericName = document.createElement("p");
    genericName.textContent = `(${unitToShow.Identity.GenericName})`;
    genericName.style.fontStyle = "italic";
    genericName.style.marginBottom = "1rem";
    detailsSection.appendChild(genericName);
  }

  const statsList = document.createElement("dl");

  // Function to add a stat if it exists
  const addStat = (label, value, parentElement) => {
    if (
      value === undefined ||
      value === null ||
      (typeof value === "string" && value === "")
    )
      return;

    const dt = document.createElement("dt");
    dt.textContent = label;
    // Ensure parentElement is valid before appending
    if (!parentElement) {
      console.error(
        "addStat called without a valid parentElement for label:",
        label,
      );
      return;
    }
    parentElement.appendChild(dt);

    const dd = document.createElement("dd");
    if (value instanceof HTMLElement) {
      dd.appendChild(value);
    } else {
      renderStats(value, dd); // Use the recursive rendering function for other types
    }
    parentElement.appendChild(dd);
  };

  // Populate with Identity details
  if (unitToShow.Identity) {
    addStat("Civilisation", unitToShow.Identity.Civ, statsList);
    addStat("Rank", unitToShow.Identity.Rank, statsList);
    addStat(
      "Classes",
      unitToShow.Identity.VisibleClasses || unitToShow.Identity.Classes,
      statsList,
    );
    addStat("Phenotype", unitToShow.Identity.Phenotype, statsList);

    if (unitToShow.Identity.Icon) {
      const iconImg = document.createElement("img");
      iconImg.src = `static/${unitToShow.Identity.Icon}`;
      iconImg.alt = (unitToShow.Identity.SpecificName || "Unit") + " icon"; // Added fallback for alt text
      iconImg.style.maxWidth = "100px";
      iconImg.style.maxHeight = "100px";
      addStat("Icon", iconImg, statsList); // Use addStat for the icon
    }
  }

  // Add other top-level stats (example, you'll need to adjust based on actual data structure)
  // This iterates through all keys of the unit object, excluding 'Identity' which is handled above
  // and internal properties like 'hasMultipleVariants' or 'code' (if it's just an internal helper)
  for (const key in unitToShow) {
    if (
      unitToShow.hasOwnProperty(key) &&
      key !== "Identity" &&
      key !== "hasMultipleVariants" &&
      key !== "code"
    ) {
      const value = unitToShow[key];
      const displayKey = key.charAt(0).toUpperCase() + key.slice(1);
      addStat(displayKey, value, statsList); // Pass statsList as parent
    }
  }

  detailsSection.appendChild(statsList);
  detailsSection.style.display = "block"; // Show the details section
}

// New recursive function to render stats, handling objects and arrays
function renderStats(data, parentElement) {
  if (typeof data === "object" && data !== null) {
    if (Array.isArray(data)) {
      // Render array as a list
      const ul = document.createElement("ul");
      ul.className = "stat-array"; // Add class for styling
      data.forEach((item) => {
        const li = document.createElement("li");
        renderStats(item, li); // Recursively render array items
        ul.appendChild(li);
      });
      parentElement.appendChild(ul);
    } else {
      // Render object as a nested definition list
      const dl = document.createElement("dl");
      dl.className = "stat-object-nested"; // Add class for styling
      for (const key in data) {
        if (data.hasOwnProperty(key)) {
          const dt = document.createElement("dt");
          dt.textContent = key.charAt(0).toUpperCase() + key.slice(1); // Capitalize key
          dl.appendChild(dt);
          const dd = document.createElement("dd");
          renderStats(data[key], dd); // Recursively render object values
          dl.appendChild(dd);
        }
      }
      parentElement.appendChild(dl);
    }
  } else {
    // Render primitive values as text content
    parentElement.textContent = data;
  }
}




function handleCivChange() {
  console.log("Civ change handler called");

  const unitsList = document.getElementById("units-list");
  const spinner = document.getElementById("loading-spinner");

  const civCode = localStorage.getItem("selectedCiv");
  if (civCode) {


    if (spinner) {
      spinner.style.display = "block"; // Show spinner
    }
    unitsList.innerHTML = ""; // Clear previous units
    // Ensure the spinner is the only child if unitsList was cleared and spinner exists
    if (spinner && !unitsList.contains(spinner)) {
      unitsList.appendChild(spinner);
    }
    fetch(`${api}/civilisations/${civCode}/units`)
      .then((res) => res.json())
      .then((units) => {
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

        units.forEach((originalUnit) => {
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

        // Store the full units data globally or pass it appropriately
        // For simplicity, let's assume a global or accessible scope for 'units'
        // This might need a more robust state management in a larger app
        window.allUnitsForCiv = units; // Store all units for the current civ

        // Group units by type using uniqueDisplayUnits
        const groups = {
          infantry: [],
          cavalry: [],
          champion: [],
          hero: [],
          ship: [],
          siege: [],
          support: [],
          other: [],
        };
        uniqueDisplayUnits.forEach((unit) => {
          // MODIFIED: Use uniqueDisplayUnits here
          const code = unit.code.toLowerCase();
          if (code.includes("infantry")) groups.infantry.push(unit);
          else if (code.includes("cavalry")) groups.cavalry.push(unit);
          else if (code.includes("champion")) groups.champion.push(unit);
          else if (code.includes("hero")) groups.hero.push(unit);
          else if (code.startsWith("ship")) groups.ship.push(unit);
          else if (code.startsWith("siege")) groups.siege.push(unit);
          else if (code.startsWith("support")) groups.support.push(unit);
          else groups.other.push(unit);
        });

        // Hide the spinner now that data is processed and before adding new elements
        if (spinner) {
          spinner.style.display = "none";
        }

        const groupOrder = [
          { key: "infantry", label: "Infantry" },
          { key: "cavalry", label: "Cavalry" },
          { key: "champion", label: "Champions" },
          { key: "hero", label: "Heroes" },
          { key: "ship", label: "Ships" },
          { key: "siege", label: "Siege" },
          { key: "support", label: "Support" },
          { key: "other", label: "Other" },
        ];
        groupOrder.forEach((group) => {
          if (groups[group.key].length) {
            const groupTitle = document.createElement("h3");
            groupTitle.textContent = group.label;
            unitsList.appendChild(groupTitle);
            const groupDiv = document.createElement("div");
            groupDiv.className = "units-group";
            groupDiv.style.display = "flex";
            groupDiv.style.flexWrap = "wrap";
            groupDiv.style.gap = "1rem";
            groups[group.key].forEach((unit) => {
              const card = document.createElement("div");
              card.className = "unit-card";
              card.tabIndex = 0;

              // Add unit icon
              if (unit.Identity && unit.Identity.Icon) {
                const icon = document.createElement("img");
                icon.src = `static/${unit.Identity.Icon}`;
                icon.alt = unit.Identity.SpecificName || "Unit icon";
                icon.style.width = "128px"; // Adjust as needed
                icon.style.height = "128px"; // Adjust as needed
                icon.style.marginBottom = "0.5rem";
                card.appendChild(icon);
              }

              // Add SpecificName
              if (unit.Identity && unit.Identity.SpecificName) {
                const specificNameElement = document.createElement("div");
                specificNameElement.textContent =
                  unit.Identity.SpecificName;
                specificNameElement.style.fontWeight = "bold";

                if (unit.hasMultipleVariants) {
                  const multiIcon = document.createElement("span");
                  multiIcon.textContent = " ➕"; // Using a plus icon, with a leading space
                  multiIcon.title =
                    "Multiple variants or upgrades exist for this unit type";
                  multiIcon.style.fontSize = "0.9em"; // Make icon slightly smaller
                  multiIcon.style.fontWeight = "normal"; // Keep icon normal weight
                  specificNameElement.appendChild(multiIcon);
                }
                card.appendChild(specificNameElement);
              }

              // Add GenericName
              if (
                unit.Identity &&
                (unit.Identity.GenericName ||
                  unit.Identity.VisibleClasses)
              ) {
                const genericNameElement = document.createElement("div");
                genericNameElement.textContent =
                  unit.Identity.VisibleClasses ??
                  unit.Identity.GenericName;
                genericNameElement.style.fontSize = "0.8em";
                genericNameElement.style.color = "#666";
                card.appendChild(genericNameElement);
              }

              card.addEventListener("click", () => {
                displayUnitDetails(unit, window.allUnitsForCiv);
              });
              card.addEventListener("keydown", (e) => {
                if (e.key === "Enter" || e.key === " ") {
                  displayUnitDetails(unit, window.allUnitsForCiv);
                }
              });

              groupDiv.appendChild(card);
            });
            unitsList.appendChild(groupDiv);
          }
        });
      })
      .catch((err) => {
        console.error("Failed to load units:", err);
        if (spinner) {
          spinner.style.display = "none"; // Hide spinner on error
        }
      });

  } else {
    unitsList.innerHTML = ""; // Clear units-list section
    if (spinner) {
      spinner.style.display = "none"; // Hide spinner if no civ is selected
    }
  }
}