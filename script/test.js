function run(input){
    // console.log(input.item)
    // console.log(input.all())
    // console.log(input.first())
    // console.log(input.last())


    return input.last()['p1']

    // test 1: read input_params
    // return JSON.stringify(input.item)
    // return JSON.stringify(input.all())
    // return input.first()

    // test 2: read files
    const fs = require('fs');
    dir_test = './'
    const files = fs.readdirSync(dir_test);
    return JSON.stringify(files)
}
