import { defineComponent } from "vue";
import CheckoutRepositoryForm from "@/pipelines/steps/components/CheckoutRepositoryForm.vue"

const stepTypes = [{
    id: 0,
    name: 'Checkout repository',
    type: 'Checkout repository',
    props: {

    },
    conditions: [
        [
            'Step Type',
            'in',
            [
                '0',
            ],
        ],
    ],
},
{
    id: 1,
    name: 'Load Training Dataset',
    props: {

    }
},
{
    id: 2,
    name: 'Train Model',
    props: {

    }
}];

export {
    stepTypes
}