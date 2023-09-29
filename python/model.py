from dataclasses import dataclass


@dataclass
class Album:
    uid: str

    @classmethod
    def from_dict(cls, data):
        return cls(
            uid=data.get("UID")
        )


@dataclass
class Photo:
    uid: str

    @classmethod
    def from_dict(cls, data):
        return cls(
            uid=data.get("PhotoUID")
        )
