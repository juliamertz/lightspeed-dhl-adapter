use std::str::FromStr;

use uuid::Uuid;

pub trait ToUuid {
    fn to_uuid(&self) -> Result<Uuid, uuid::Error>;
}

impl<T> ToUuid for T
where
    T: AsRef<str>,
{
    fn to_uuid(&self) -> Result<Uuid, uuid::Error> {
        Uuid::from_str(self.as_ref())
    }
}
