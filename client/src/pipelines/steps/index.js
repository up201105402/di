import { defineComponent } from "vue";

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
    component: defineComponent({
        name: 'CheckoutRepositoryForm',
        setup() {
          useAuth();
        }
      })
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