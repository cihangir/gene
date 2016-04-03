var iz, processRequest;
iz = require('iz');
processRequest = require('./_request.js');
module.exports = (function(o) {
  return {
    ByFacebookIDs: function(data, callback) {
      rules = {
        v: iz(i.v).required().minLength(1)
      };
      areRules = are(rules);
      for (var i = 0; i < data.length; i++){
        if (!areRules.validFor(data[i])){
          return callback(areRules.getInvalidFields());
        }
      }
      return processRequest(o.baseUrl + '/account/byfacebookids', data, callback)
    },
    ByIDs: function(data, callback) {
      rules = {
        v: iz(null).required().minLength(1)
      };
      areRules = are(rules);
      for (var i = 0; i < data.length; i++){
        if (!areRules.validFor(data[i])){
          return callback(areRules.getInvalidFields());
        }
      }
      return processRequest(o.baseUrl + '/account/byids', data, callback)
    },
    Create: function(data, callback) {
      rules = {
        deletedAt: iz(a.deletedAt).required().date(),
        facebookSecretToken: iz(a.facebookSecretToken).required().minLength(1),
        id: iz(null).required().minLength(1),
        statusConstant: iz(a.statusConstant).required().inArray(["registered","disabled","spam",]),
        createdAt: iz(a.createdAt).required().date(),
        emailStatusConstant: iz(a.emailStatusConstant).required().inArray(["verified","notVerified",]),
        facebookAccessToken: iz(a.facebookAccessToken).required().minLength(1),
        facebookID: iz(a.facebookID).required().minLength(1),
        profileID: iz(null).required().minLength(1),
        updatedAt: iz(a.updatedAt).required().date()
      };
      areRules = are(rules);
      if (!areRules.validFor(data)){
        return callback(areRules.getInvalidFields());
      }
      return processRequest(o.baseUrl + '/account/create', data, callback)
    },
    Delete: function(data, callback) {
      return processRequest(o.baseUrl + '/account/delete', data, callback)
    },
    One: function(data, callback) {
      return processRequest(o.baseUrl + '/account/one', data, callback)
    },
    Update: function(data, callback) {
      rules = {
        emailStatusConstant: iz(a.emailStatusConstant).required().inArray(["verified","notVerified",]),
        facebookAccessToken: iz(a.facebookAccessToken).required().minLength(1),
        facebookID: iz(a.facebookID).required().minLength(1),
        profileID: iz(null).required().minLength(1),
        updatedAt: iz(a.updatedAt).required().date(),
        createdAt: iz(a.createdAt).required().date(),
        facebookSecretToken: iz(a.facebookSecretToken).required().minLength(1),
        id: iz(null).required().minLength(1),
        statusConstant: iz(a.statusConstant).required().inArray(["registered","disabled","spam",]),
        deletedAt: iz(a.deletedAt).required().date()
      };
      areRules = are(rules);
      if (!areRules.validFor(data)){
        return callback(areRules.getInvalidFields());
      }
      return processRequest(o.baseUrl + '/account/update', data, callback)
    },
  }
})(o);