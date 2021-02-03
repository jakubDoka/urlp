# url parser 

Does similar thing to json build in package but it does not use bytes but url.Values. This can come handy if you are sending lot of numerical data or booleans through url. It also removes `values["something"][0]` that is usually spread all over the place. Biggest advantage is that struct is argument assertion by it self though you can make fields optional with annotations.