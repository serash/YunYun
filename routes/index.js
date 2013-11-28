
/*
 * GET home page.
 */

exports.index = function(req, res){
  res.render('index', { title: 'YunYun' });
};

exports.yunyun = function(req, res){
  res.render('index', { title: 'YunYun' });
};

exports.review = function(req, res){
  res.render('review', { title: 'Review' });
};

exports.learn = function(req, res){
  res.render('learn', { title: 'Learn' });
};

exports.stats = function(req, res){
  res.render('stats', { title: 'Stats' });
};

exports.user = function(db) {
    return function(req, res) {
        var collection = db.get('usercollection');
        collection.find({},{},function(e,docs){
            res.render('user', {
                "userlist" : docs,
                 title: 'User' 
            });
        });
    };
};