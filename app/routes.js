//SPDX-License-Identifier: Apache-2.0

var controller = require('./controller.js');

module.exports = function(app){

  // The cum-group area
 app.get('/get_all_groups', function(req, res){
    controller.get_all_groups(req, res);
  });
  app.get('/add_group/:newGroup', function(req, res){
    controller.add_group(req, res);
  });

  // The cum-rec area
  app.get('/add_user/:user', function(req, res){
    controller.add_user(req, res);
  });
  app.get('/query_all_users', function(req, res){
    controller.query_all_users(req, res);
  });
  app.get('/generate_set_for_group/:generator', function(req, res){
    controller.generate_set_for_group(req, res);
  });
  app.get('/get_user_record/:id', function(req, res){
      controller.get_user_record(req, res);
  });
  app.get('/prepare_for_delivery/:exam', function(req, res){
      controller.prepare_for_delivery(req, res);
  });
  app.get('/delivery_item/:delicase', function(req, res){
      controller.delivery_item(req, res);
  });

}
