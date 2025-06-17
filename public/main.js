const api = "";

function unitComponent() {
    return {
        savedCiv: localStorage.getItem('selectedCiv'),
        unitsClass: {},
        selectedUnit: null,
        unitDetailsArr: [],
        unitDetails: null,
        displayVersion: 0,
        ressources: [],

        async init() {
            this.$watch('savedCiv', async () => {
                await this.loadUnits();
            });

            this.$watch('selectedUnit', () => {
                this.fetchUnitDetails();
            });

            if (this.savedCiv) {
                await this.loadUnits();
            }
        },

        async loadUnits() {
            const response = await fetch(`/civilisations/${this.savedCiv}/units`);
            this.unitsClass = await response.json();
        },

        async fetchUnitDetails() {
            const res = await fetch(`/civilisations/${this.savedCiv}/units/${this.selectedUnit}`);
            this.unitDetailsArr = await res.json();
            this.displayVersion = 0;
            this.unitDetails = this.unitDetailsArr[0];
            this.setRessources();
        },

        setDisplayVersion(version) {
            this.displayVersion = version;
            this.unitDetails = this.unitDetailsArr[version];
            this.setRessources();
        },

        setRessources() {
            const cost = this.unitDetails.Cost || {};
            const resources = cost.Resources || {};

            this.ressources = [
                { img: '/static/assets/icons/population.png', value: cost.Population || 0 },
                { img: '/static/assets/icons/food.png', value: resources.food || 0 },
                { img: '/static/assets/icons/wood.png', value: resources.wood || 0 },
                { img: '/static/assets/icons/metal.png', value: resources.metal || 0 },
            ];
        }
    }
}
