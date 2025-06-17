const api = "";

function civUnitBrowser() {
    return {
        selectedCivilisation: localStorage.getItem('selectedCiv'),
        unitGroups: {},
        activeUnitCode: null,
        unitVariants: [],
        currentUnit: null,
        currentVariantIndex: 0,
        resourceCosts: [],

        async initialize() {
            this.$watch('selectedCivilisation', async () => {
                await this.loadUnitsForCivilisation();
            });

            this.$watch('activeUnitCode', () => {
                this.loadUnitVariantDetails();
            });

            if (this.selectedCivilisation) {
                await this.loadUnitsForCivilisation();
            }
        },

        async loadUnitsForCivilisation() {
            const response = await fetch(`/civilisations/${this.selectedCivilisation}/units`);
            this.unitGroups = await response.json();
        },

        async loadUnitVariantDetails() {
            const res = await fetch(`/civilisations/${this.selectedCivilisation}/units/${this.activeUnitCode}`);
            this.unitVariants = await res.json();
            this.currentVariantIndex = 0;
            this.currentUnit = this.unitVariants[0];
            this.updateResourceCosts();
        },

        switchUnitVariant(variantIdx) {
            this.currentVariantIndex = variantIdx;
            this.currentUnit = this.unitVariants[variantIdx];
            this.updateResourceCosts();
        },

        updateResourceCosts() {
            const cost = this.currentUnit.Cost || {};
            const resources = cost.Resources || {};

            this.resourceCosts = [
                { img: '/static/assets/icons/population.png', value: cost.Population || 0 },
                { img: '/static/assets/icons/food.png', value: resources.food || 0 },
                { img: '/static/assets/icons/wood.png', value: resources.wood || 0 },
                { img: '/static/assets/icons/metal.png', value: resources.metal || 0 },
            ];
        }
    }
}
