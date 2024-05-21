# coding=utf-8
# *** WARNING: this file was generated by pulumi-language-python. ***
# *** Do not edit by hand unless you're certain you know what you are doing! ***

import copy
import warnings
import pulumi
import pulumi.runtime
from typing import Any, Mapping, Optional, Sequence, Union, overload
from . import _utilities

__all__ = [
    'DnsRecord',
]

@pulumi.output_type
class DnsRecord(dict):
    def __init__(__self__, *,
                 hostname: str,
                 type: str,
                 value: str):
        pulumi.set(__self__, "hostname", hostname)
        pulumi.set(__self__, "type", type)
        pulumi.set(__self__, "value", value)

    @property
    @pulumi.getter
    def hostname(self) -> str:
        return pulumi.get(self, "hostname")

    @property
    @pulumi.getter
    def type(self) -> str:
        return pulumi.get(self, "type")

    @property
    @pulumi.getter
    def value(self) -> str:
        return pulumi.get(self, "value")


